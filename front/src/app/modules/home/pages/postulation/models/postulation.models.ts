export interface SocioEconomicSheet {
  sections: SocioSection[]
}

interface SocioSection {
  name: string;
  type_section: string;
  attributes: SocioAttribute[]
}

export interface SocioAttribute {
  name: string;
  description: string;
  value: string;
  type_attribute: string;
  options?: string;
  order: number
}

