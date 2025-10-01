export interface Response<T = any> {
  error: boolean;
  msg: string;
  data: T;
  code: number;
  type: string;
  metadata: Metadata;
}

export interface ItemBasic {
  name: string;
  value: string;
}

export interface Metadata {
  links: Links;
  total: number;
  offset: number;
  limit: number;
  page: number;
  total_pages: number;
}

export interface Links {
  next: string;
  prev: string;
}
