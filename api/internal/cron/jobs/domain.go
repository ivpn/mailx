package jobs

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/repository"
	"ivpn.net/email/api/internal/service"
)

const domainVerifyBatchSize = 100

// VerifyDomainsJob checks ownership and DNS records for all domains.
// Fields owner_verified_at, mx_verified_at and send_verified_at are set to NULL
// when the respective check fails. Domains are processed in batches of 100 with
// a 200ms sleep between each domain to avoid saturating DNS resolvers or the DB.
func VerifyDomainsJob(cfg config.Config, db *gorm.DB) {
	log.Println("VerifyDomainsJob: starting")

	repo := &repository.Database{Client: db}
	svc := service.New(cfg, repo, nil)
	ctx := context.Background()

	offset := 0
	total := 0

	for {
		var batch []model.Domain
		if err := db.Order("id").Limit(domainVerifyBatchSize).Offset(offset).Find(&batch).Error; err != nil {
			log.Printf("VerifyDomainsJob: error fetching domains at offset %d: %s", offset, err)
			return
		}

		if len(batch) == 0 {
			break
		}

		for _, domain := range batch {
			// ownership check
			if err := svc.VerifyDomainOwner(ctx, domain.Name, domain.UserID); err != nil {
				log.Printf("VerifyDomainsJob: ownership check failed for domain %s: %s", domain.Name, err)
				if dbErr := db.Model(&model.Domain{}).Where("id = ?", domain.ID).Updates(map[string]any{
					"owner_verified_at": nil,
				}).Error; dbErr != nil {
					log.Printf("VerifyDomainsJob: error nulling owner_verified_at for domain %s: %s", domain.Name, dbErr)
				}
			}

			// MX records check
			if err := svc.VerifyDomainMX(ctx, domain.Name, domain.UserID); err != nil {
				log.Printf("VerifyDomainsJob: MX check failed for domain %s: %s", domain.Name, err)
				if dbErr := db.Model(&model.Domain{}).Where("id = ?", domain.ID).Updates(map[string]any{
					"mx_verified_at": nil,
				}).Error; dbErr != nil {
					log.Printf("VerifyDomainsJob: error nulling mx_verified_at for domain %s: %s", domain.Name, dbErr)
				}
			}

			// Send records check (SPF, DKIM, DMARC)
			if err := svc.VerifyDomainSend(ctx, domain.Name, domain.UserID); err != nil {
				log.Printf("VerifyDomainsJob: send records check failed for domain %s: %s", domain.Name, err)
				if dbErr := db.Model(&model.Domain{}).Where("id = ?", domain.ID).Updates(map[string]any{
					"send_verified_at": nil,
				}).Error; dbErr != nil {
					log.Printf("VerifyDomainsJob: error nulling send_verified_at for domain %s: %s", domain.Name, dbErr)
				}
			}

			time.Sleep(200 * time.Millisecond)
		}

		total += len(batch)
		offset += len(batch)
	}

	log.Printf("VerifyDomainsJob: completed, processed %d domains", total)
}
