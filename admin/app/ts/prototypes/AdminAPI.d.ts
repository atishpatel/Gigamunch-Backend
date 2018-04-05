interface ErrorOnlyResp {
    error: Error
}
interface TokenOnlyReq {
    token: string
}
interface TokenOnlyResp {
    error: Error
    token: string
}
interface GetLogReq {
    id: number
}
interface GetLogResp {
    error: Error
    log: Log
}
interface GetLogsReq {
    start: number
    limit: number
}
interface GetLogsResp {
    error: Error
    logs: Log[]
}
interface Sublog {
    date: string
    sub_email: string
    created_datetime: string
    skip: boolean
    servings: number
    amount: number
    amount_paid: number
    paid: boolean
    paid_datetime: string
    payment_method_token: string
    transaction_id: string
    free: boolean
    discount_amount: number
    discount_percent: number
    customer_id: string
    refunded: boolean
}
interface Subscriber {
    email: string
    date: string
    name: string
    first_name: string
    last_name: string
    address: Address.Address
    customer_id: string
    subscription_ids: string[]
    first_payment_date: string
    is_subscribed: boolean
    subscription_date: string
    unsubscribed_date: string
    first_box_date: string
    servings: number
    vegetarian_servings: number
    delivery_time: number
    subscription_day: string
    weekly_amount: number
    payment_method_token: string
    reference: string
    phone_number: string
    delivery_tips: string
    bag_reminder_sms: boolean
    // gift
    num_gift_dinners: number
    reference_email: string
    gift_reveal_date: string
    // stats
    referral_page_opens: number
    referred_page_opens: number
    gift_page_opens: number
    gifted_page_opens: number
}
interface GetUnpaidSublogsReq {
    limit: number
}
interface GetUnpaidSublogsResp {
    error: Error
    sublogs: Sublog[]
}
interface ProcessSublogsReq {
    email: string
    date: string
}
interface ProcessSublogsResp {
    error: Error
}
interface GetAllSubscribersReq {
    date: string
}
interface GetAllSubscribersResp {
    error: Error
    subscribers: Subscriber[]
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
interface GetAllExecutionsReq {
    start: number
    limit: number
}
interface GetAllExecutionsResp {
    error: Error
    execution: Execution[]
}
interface GetExecutionReq {
   id: number
}
interface GetExecutionResp {
    error: Error
    execution: Execution
}
interface UpdateExecutionReq {
    execution: Execution
}
interface UpdateExecutionResp {
    error: Error
    execution: Execution
}
interface ExecutionStats {
    id: number
    created_datetime: string
    date: string
    location: number
    nationality: string
    country: string
    city: string
    revenue: number
    payroll: Payroll
    payroll_costs: number
    food_costs: number
    delivery_costs: number
    onboarding_costs: number
    processing_costs: number
    tax_costs: number
    packaging_costs: number
    other_costs: number
}
interface Payroll {
    name: string
    hours: number
    wage: number
    postion: string
}
interface GetAllExecutionStatsReq {
    start: number
    limit: number
}
interface GetAllExecutionStatsResp {
    error: Error
    execution_stats: ExecutionStats[]
}
interface GetExecutionStatsReq {
    id: number
}
interface GetExecutionStatsResp {
    error: Error
    execution_stats: ExecutionStats
}
interface UpdateExecutionStatsReq {
    execution_stats: Execution
}
interface UpdateExecutionStatsResp {
    error: Error
    execution_stats: Execution
}
