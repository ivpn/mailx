export interface Alias {
  id: string
  name: string
  enabled: boolean
  description: string
  recipients: string
  domain: string
  format: string
  from_name: string
  wildcard: boolean
  local_part: string
  pinned: boolean
}

export interface CustomDomain {
  id: string
  name: string
  enabled: boolean
  owner_verified_at: string | null
  mx_verified_at: string | null
  send_verified_at: string | null
}

export interface Defaults {
  domain: string
  domains: string[]
  custom_domains: CustomDomain[]
  alias_format: string
  recipient: string
  recipients: string[]
}

export interface Preferences {
  input_button: boolean
  add_description: boolean
}
