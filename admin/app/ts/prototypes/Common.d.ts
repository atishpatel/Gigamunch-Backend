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
interface Log {
  id: number
  log_name: string
  timestamp: string
  type: string
  action: string
  action_user_id: number
  action_user_email: string
  user_id: number
  user_email: string
  severity: number
  path: string
  basic_payload: BasicPayload
}
interface BasicPayload {
  title: string
  description: string
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
interface ExecutionProgress {
  head_chef: number
  photographer: number
  content_writer: number
  culture_guide: number
}
interface Email {
	dinner_non_veg_image_url: string
	dinner_veg_image_url: string
	cook_image_url: string
	landscape_image_url: string
}
interface Notifications {
  delivery_sms: string
  rating_sms: string
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
}