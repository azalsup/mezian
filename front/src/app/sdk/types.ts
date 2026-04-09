// ── Auth ──────────────────────────────────────────────────────────────────────

export interface AuthTokens {
  access_token:  string;
  refresh_token: string;
  token_type:    'Bearer';
  expires_in:    number;
}

export interface User {
  id:           number;
  phone:        string;
  email?:       string;
  is_verified:  boolean;
  display_name: string;
  avatar_url?:  string;
  address?:     string;
  city?:        string;
  postal_code?: string;
  country?:     string;
  role:         'user' | 'admin' | 'moderator' | string;
  roles?:       Role[];
  created_at:   string;
  updated_at:   string;
}

export interface AuthResponse {
  tokens: AuthTokens;
  user:   User;
}

export interface RegisterPayload {
  phone?:        string;
  email?:        string;
  password:      string;
  display_name:  string;
  address?:      string;
  city?:         string;
  postal_code?:  string;
  country?:      string;
}

// ── RBAC ──────────────────────────────────────────────────────────────────────

export interface Permission {
  id:          number;
  key:         string;   // e.g. "categories.write"
  group:       string;   // e.g. "categories"
  label_fr:    string;
  description: string;
}

export interface Role {
  id:           number;
  name:         string;
  slug:         string;
  description:  string;
  is_system:    boolean;
  permissions?: Permission[];
  created_at:   string;
  updated_at:   string;
}

export interface UserListResponse {
  data:      User[];
  total:     number;
  page:      number;
  page_size: number;
}

// ── Ads ───────────────────────────────────────────────────────────────────────

export interface Ad {
  id:                number;
  title:             string;
  price:             number | null;
  city:              string;
  category_slug:     string;
  subcategory_slug?: string;
  images:            string[];
  created_at:        string;
  badge?:            'top' | 'new' | 'urgent' | 'premium';
  description?:      string;
  seller_name?:      string;
  seller_phone?:     string;
  views?:            number;
}

export interface AdsQuery {
  q?:        string;
  cat?:      string;
  sub?:      string;
  city?:     string;
  minPrice?: number;
  maxPrice?: number;
  page?:     number;
  sort?:     'newest' | 'price_asc' | 'price_desc';
}

export interface AdsResponse {
  data:      Ad[];
  total:     number;
  page:      number;
  page_size: number;
}

export interface CreateAdPayload {
  category_id:    number;
  title:          string;
  body:           string;
  price?:         number;
  currency:       string;
  city:           string;
}

// ── Categories ────────────────────────────────────────────────────────────────

export interface AttributeDefinition {
  id:            number;
  key:           string;
  label_fr:      string;
  label_ar:      string;
  label_en:      string;
  data_type:     'string' | 'integer' | 'float' | 'boolean' | 'enum';
  unit?:         string;
  enum_values?:  string[];
  is_required:   boolean;
  is_filterable: boolean;
  sort_order:    number;
}

export interface Category {
  id:                     number;
  slug:                   string;
  name_fr:                string;
  name_ar:                string;
  name_en:                string;
  icon?:                  string;
  sort_order:             number;
  is_active:              boolean;
  featured?:              boolean;
  parent_id?:             number;
  children?:              Category[];
  attribute_definitions?: AttributeDefinition[];
}
