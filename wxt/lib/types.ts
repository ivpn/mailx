export interface Alias {
  id: string
  name: string
  enabled: boolean
  description: string
  recipients: string
}

export interface Defaults {
  domain: string
  recipient: string
  alias_format: string
  recipients: string
}
