declare namespace Common {
interface Error {
  code: number
  message: string
  detail: string
}
interface Address {
  apt: string
  street: string
  city: string
  state: string
  zip: string
  country: string
  latitude: number
  longitude: number
  full_address: string
}
interface Campaign {
  Source: string
  Medium: string
  Campaign: string
  Term: string
  Content: string
  Timestamp: string
}
interface Log {
  id: number
  log_name: string
  timestamp: string
  type: string
  action: string
  action_user_id: string
  action_user_email: string
  user_id: string
  user_email: string
  severity: number
  path: string
  basic_payload: BasicPayload
}
interface BasicPayload {
  title: string
  description: string
}
interface Activity {
  created_datetime: string
  date: string
  user_id: string
  email: string
  first_name: string
  last_name: string
  location: number
  address_changed: boolean
  address_apt: string
  address_string: string
  zip: string
  latitude: number
  longitude: number
  active: boolean
  skip: boolean
  forgiven: boolean
  servings_non_vegetarian: number
  servings_vegetarian: number
  servings_changed: boolean
  first: boolean
  amount: number
  amount_paid: number
  discount_amount: number
  discount_percent: number
  paid: boolean
  paid_datetime: string
  payment_provider: number
  transaction_id: string
  payment_method_token: string
  customer_id: string
  refunded: boolean
  refunded_amount: number
  refunded_datetime: string
  refund_transaction_id: string
  gift: boolean
  gift_from_user_id: number
  deviant: boolean
  deviant_reason: string
}
interface Discount {
	id: number
	created_datetime: string
	user_id: string
	email: string
	first_name: string
	last_name: string
	date_used: string
	discount_amount: number
	discount_percent: number
}
interface Execution {
  id: number
  date: string
  location: number
  publish: boolean
  created_datetime: string
  culture: Culture
  content: Content
  culture_cook: CultureCook
  culture_guide: CultureGuide
  dishes: Dish[]
  stickers: Sticker[]
  notifications: Notifications
  has_pork: boolean
  has_beef: boolean
  has_chicken: boolean
  email: Email
}
interface ExecutionAndActivity {
  execution: Execution
  activity: Activity
}
interface ExecutionProgressSummary {
  message: string
  is_error: boolean
}
interface ExecutionProgress {
  head_chef: number
  content_writer: number
  culture_guide: number
  summary: ExecutionProgressSummary[]
}
interface Email {
  dinner_non_veg_image_url: string
  dinner_veg_image_url: string
  cook_image_url: string
  landscape_image_url: string
  cook_face_image_url: string
}
interface Notifications {
  delivery_sms: string
  rating_sms: string
  rating_link_veg: string
  rating_link_nonveg: string
}
interface InfoBox {
  title: string
  text: string
  caption: string
  image: string
}
interface CultureGuide {
  info_boxes: InfoBox[]
  dinner_instructions: string
  main_color: string
  font_name: string
  font_style: string
  font_caps: boolean
  vegetarian_dinner_instructions: string
  font_name_post_script: string
}
interface Content {
  landscape_image_url: string
  cook_image_url: string
  hands_plate_non_veg_image_url: string
  hands_plate_veg_image_url: string
  dinner_non_veg_image_url: string
  spotify_url: string
  youtube_url: string
  font_url: string
  dinner_veg_image_url: string
  cover_image_url: string
  map_image_url: string
  cook_face_image_url: string
}
interface Culture {
  country: string
  city: string
  description: string
  nationality: string
  greeting: string
  flag_emoji: string
  description_preview: string
}
interface Sticker {
  name: string
  ingredients: string
  extra_instructions: string
  reheat_option_1: string
  reheat_option_2: string
  reheat_time_1: string
  reheat_time_2: string
  reheat_instructions_1: string
  reheat_instructions_2: string
  eating_temperature: string
  reheat_option_1_preferred: boolean
  number: number
  color: string
  is_for_non_vegetarian: boolean
  is_for_vegetarian: boolean
}
interface Dish {
  number: number
  color: string
  name: string
  description: string
  ingredients: string
  is_for_vegetarian: boolean
  is_for_non_vegetarian: boolean
  is_on_main_plate: boolean
  image_url: string
  description_preview: string
  container_size: string
}
interface QandA {
  question: string
  answer: string
}
interface CultureCook {
  first_name: string
  last_name: string
  story: string
  story_preview: string
  q_and_a: QandA[]
}
interface Subscriber {
  created_datetime: string
  sign_up_datetime: string
  id: string
  auth_id: string
  location: number
  photo_url: string
  // Pref
  email_prefs: EmailPref[]
  phone_prefs: PhonePref[]
  // Account
  payment_provider: number
  payment_customer_id: string
  payment_method_token: string
  active: boolean
  activate_datetime: string
  deactivated_datetime: string
  address: Address
  delivery_notes: string
  // Plan
  servings_non_vegetarian: number
  servings_vegetarian: number
  plan_interval: number
  plan_weekday: string
  interval_start_point: string
  amount: number
  food_pref: FoodPref
  // Gift
  num_gift_dinners: number
  gift_reveal_datetime: string
  // Marketing
  referral_page_opens: number
  referred_page_opens: number
  referrer_user_id: number
  reference_email: string
  reference_text: string
  campaigns: Campaign[]
}
interface FoodPref {
  no_pork: boolean
  no_beef: boolean
}
interface EmailPref {
  default: boolean
  first_name: string
  last_name: string
  email: string
}
interface PhonePref {
  number: string
  raw_number: string
  disable_bag_reminder: boolean
  disable_delivered: boolean
  disable_review: boolean
}
interface ErrorOnlyResp {
  error: Error
}
}