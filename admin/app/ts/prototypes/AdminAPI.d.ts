interface ErrorOnlyResp {
    error: Error
}
interface MakeAdminReq {
    email: string
}
interface Activity {
    created_datetime: string
    date: string
    user_id: number
    email: string
    first_name: string
    last_name: string
    location: number
    // Address
    address_changed: boolean
    address_apt: string
    address_string: string
    zip: string
    latitude: number
    longitude: number
    // Detail
    active: boolean
    skip: boolean
    // Bag detail
    servings: number
    vegetarian_servings: boolean
    servings_changed: number
    first: boolean
    // Payment
    amount: number
    amount_paid: number
    discount_amount: number
    discount_percent: number
    paid: boolean
    paid_datetime: string
    transaction_id: string
    payment_method_token: string
    customer_id: string
    // Refund
    refunded: boolean
    refunded_amount: number
    refunded_datetime: string
    refund_transaction_id: string
    payment_provider: number
    forgiven: boolean
    // Gift
    gift: boolean
    gift_from_user_id: number
    // Deviant
    deviant: boolean
    deviant_reason: string
}
interface SkipActivityReq {
    email: string
    date: string
}
interface SkipActivityResp {
    error: Error
}
interface UnskipActivityReq {
    email: string
    date: string
}
interface UnskipActivityResp {
    error: Error
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
interface GetLogsByEmailReq {
    start: number
    limit: number
    email: string
}
interface GetLogsByEmailResp {
    error: Error
    logs: Log[]
}
interface Sublog {
    date: string
    sub_email: string
    created_datetime: string
    skip: boolean
    servings: number
    veg_servings: number
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
    address: Address
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
interface SendCustomerSMSReq {
    emails: string[]
	message: string
}
interface RefundAndSkipSublogReq {
	email: string
    date: string
}
interface GetHasSubscribedReq {
    date: string
}
interface GetSubscriberSublogsReq {
    email: string
}
interface GetSubscriberSublogsResp {
    error: Error
    sublogs: Sublog[]
}
interface GetHasSubscribedResp {
    error: Error
    subscribers: Subscriber[]
}
interface GetSubscriberReq {
    email: string
}
interface GetSubscriberResp {
    error: Error
    subscriber: Subscriber
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
interface GetExecutionsReq {
    start: number
    limit: number
}
interface GetExecutionsResp {
    error: Error
    executions: Execution[]
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
interface Delivery {
    date: string
    driver_name: string
    driver_email: string
    sub_email: string
    order: number
    success: boolean
    fail: boolean
    sub_name: string
    phone_number: string
    address: Address
    delivery_notes: string
    servings: number
    vegetarian: boolean
    first: boolean
}
interface GetDeliveriesReq {
    date: string
    driver_email: string
}
interface GetDeliveriesResp {
    error: Error
    deliveries: Delivery[]
}
interface UpdateDeliveriesReq {
    deliveries: Delivery[]
}
interface UpdateDeliveriesResp {
    error: Error
}
