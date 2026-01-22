export interface Alias {
  id: string
  name: string
  enabled: boolean
  description: string
  recipients: string
  domain: string
  format: string
  from_name: string
  catch_all: boolean
}

export interface Defaults {
  domain: string
  domains: string[]
  alias_format: string
  recipient: string
  recipients: string[]
}

export interface Preferences {
  input_button: boolean
}
