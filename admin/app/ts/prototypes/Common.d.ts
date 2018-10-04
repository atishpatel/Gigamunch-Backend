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
    // Info
    culture: Culture
    content: Content
    culture_cook: CultureCook
    dishes: Dish[]
    // Diet
    has_pork: boolean
    has_beef: boolean
    has_chicken: boolean
    has_weird_meat: boolean
    has_fish: boolean
    has_other_seafood: boolean
}
interface Content {
    hero_image_url: string
    cook_image_url: string
    hands_plate_image_url: string
    dinner_image_url: string
    spotify_url: string
    youtube_url: string
}
interface Culture {
    country: string
    city: string
    description: string
    nationality: string
    greeting: string
    flag_emoji: string
}
interface Dish {
    number: number
    color: string
    name: string
    description: string
    ingredients: string
    is_for_vegetarian: boolean
    is_for_non_vegetarian: boolean
}
interface CultureCook {
    first_name: string
    last_name: string
    story: string
}
}