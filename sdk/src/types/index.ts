/**
 * Types TypeScript miroirs des structs Go du backend Mezian.
 * Chaque interface correspond exactement aux champs JSON retournés par l'API.
 */

// ---------------------------------------------------------------------------
// Utilitaires
// ---------------------------------------------------------------------------

/** Standard API error response */
export interface ApiError {
  error: string;
  details?: Record<string, string>;
}

/** Generic paginated response */
export interface Paginated<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

// ---------------------------------------------------------------------------
// User & Auth
// ---------------------------------------------------------------------------

export type UserRole = "user" | "admin";

export interface User {
  id: number;
  phone: string;
  email?: string;
  is_verified: boolean;
  display_name: string;
  avatar_url?: string;
  city?: string;
  role: UserRole;
  created_at: string; // ISO 8601
  updated_at: string;
}

/** Tokens returned after successful auth */
export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  token_type: "Bearer";
  expires_in: number; // secondes
}

/** Complete login/register response */
export interface AuthResponse {
  tokens: AuthTokens;
  user: User;
}

// Auth request bodies
export interface SendOtpRequest {
  phone: string;
  /** Canal d'envoi. Défaut: sms */
  channel?: "sms" | "whatsapp";
  /** Objectif du code. Défaut: login */
  purpose?: "login" | "signup" | "phone_change";
}

export interface VerifyOtpRequest {
  phone: string;
  code: string;
  purpose?: "login" | "signup" | "phone_change";
}

export interface LoginRequest {
  /** E.164 phone number or email */
  identifier: string;
  password: string;
}

export interface RegisterRequest {
  phone: string;
  email?: string;
  password?: string;
  display_name: string;
}

export interface UpdateMeRequest {
  display_name?: string;
  email?: string;
  city?: string;
  avatar_url?: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

// ---------------------------------------------------------------------------
// Categories & Attributes
// ---------------------------------------------------------------------------

export type AttributeDataType = "integer" | "float" | "string" | "boolean" | "enum";

export interface AttributeDefinition {
  id: number;
  category_id: number;
  key: string;
  label_fr: string;
  label_ar: string;
  data_type: AttributeDataType;
  unit?: string;
  /** JSON array de valeurs possibles, ex: '["Essence","Diesel"]' */
  enum_values?: string;
  is_required: boolean;
  is_filterable: boolean;
  sort_order: number;
}

export interface Category {
  id: number;
  slug: string;
  name_fr: string;
  name_ar: string;
  icon?: string;
  parent_id?: number;
  sort_order: number;
  is_active: boolean;
  children?: Category[];
  attribute_definitions?: AttributeDefinition[];
}

// ---------------------------------------------------------------------------
// Ads
// ---------------------------------------------------------------------------

export type AdStatus = "draft" | "active" | "sold" | "expired" | "deleted";
export type Currency = "MAD" | "EUR" | "USD";

export interface AdAttribute {
  id: number;
  ad_id: number;
  key: string;
  value: string;
}

export type MediaType = "image" | "youtube";

export interface Media {
  id: number;
  ad_id: number;
  type: MediaType;
  url: string;
  thumb_url?: string;
  sort_order: number;
  is_cover: boolean;
  created_at: string;
}

export interface Ad {
  id: number;
  user_id: number;
  category_id: number;
  shop_id?: number;
  slug: string;
  title: string;
  /** Body in Markdown */
  body: string;
  price?: number;
  currency: Currency;
  city: string;
  district?: string;
  status: AdStatus;
  is_boosted: boolean;
  boosted_until?: string;
  view_count: number;
  created_at: string;
  updated_at: string;
  // Relations
  user?: Pick<User, "id" | "display_name" | "phone" | "avatar_url">;
  category?: Pick<Category, "id" | "slug" | "name_fr" | "name_ar">;
  shop?: Pick<Shop, "id" | "slug" | "name" | "logo_url">;
  media?: Media[];
  attributes?: AdAttribute[];
}

// Ad request bodies
export interface CreateAdRequest {
  category_id: number;
  title: string;
  /** Markdown */
  body: string;
  price?: number;
  currency?: Currency;
  city: string;
  district?: string;
  status?: AdStatus;
  attributes?: Array<{ key: string; value: string }>;
}

export interface UpdateAdRequest extends Partial<CreateAdRequest> {}

export interface AdFilters {
  category_id?: number;
  city?: string;
  min_price?: number;
  max_price?: number;
  search?: string;
  status?: AdStatus;
  shop_id?: number;
  sort?: "newest" | "oldest" | "price_asc" | "price_desc" | "views";
  page?: number;
  limit?: number;
}

// ---------------------------------------------------------------------------
// Shops
// ---------------------------------------------------------------------------

export type ShopPlan = "starter" | "pro" | "premium";

export interface Shop {
  id: number;
  user_id: number;
  slug: string;
  name: string;
  /** Description in Markdown */
  description?: string;
  logo_url?: string;
  cover_url?: string;
  phone: string;
  city: string;
  plan: ShopPlan;
  plan_expires?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  user?: Pick<User, "id" | "display_name">;
}

export interface CreateShopRequest {
  name: string;
  description?: string;
  phone: string;
  city: string;
}

export interface UpdateShopRequest extends Partial<CreateShopRequest> {
  logo_url?: string;
  cover_url?: string;
}
