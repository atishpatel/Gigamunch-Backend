interface ErrorOnlyResp {
    error: Common.Error
}
interface SetAdminReq {
    email: string
    active: boolean
}
interface ActivateSubscriberReq {
    email: string
    // optional
    first_bag_date: string
}
interface DeactivateSubscriberReq {
    email: string
}
interface ReplaceSubscriberEmailReq {
    old_email: string
    new_email: string
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
    error: Common.Error
}
interface UnskipActivityReq {
    email: string
    date: string
}
interface UnskipActivityResp {
    error: Common.Error
}
interface RefundActivityReq {
    email: string
    date: string
    amount: number
    percent: number
}
interface RefundAndSkipActivityReq {
    email: string
    date: string
    amount: number
    percent: number
}
interface GetLogReq {
    id: number
}
interface GetLogResp {
    error: Common.Error
    log: Common.Log
}
interface GetLogsReq {
    start: number
    limit: number
}
interface GetLogsResp {
    error: Common.Error
    logs: Common.Log[]
}
interface GetLogsByEmailReq {
    start: number
    limit: number
    email: string
}
interface GetLogsByEmailResp {
    error: Common.Error
    logs: Common.Log[]
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
    refunded_amount: number
}
interface Subscriber {
    email: string
    date: string
    name: string
    first_name: string
    last_name: string
    address: Common.Address
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
    error: Common.Error
    sublogs: Sublog[]
}
interface ProcessSublogsReq {
    email: string
    date: string
}
interface ProcessSublogsResp {
    error: Common.Error
}
interface SendCustomerSMSReq {
    emails: string[]
	message: string
}
interface GetHasSubscribedReq {
    date: string
}
interface GetSubscriberSublogsReq {
    email: string
}
interface GetSubscriberSublogsResp {
    error: Common.Error
    sublogs: Sublog[]
}
interface GetHasSubscribedResp {
    error: Common.Error
    subscribers: Subscriber[]
}
interface GetSubscriberReq {
    email: string
}
interface GetSubscriberResp {
    error: Common.Error
    subscriber: Subscriber
}
interface GetExecutionsReq {
    start: number
    limit: number
}
interface GetExecutionsResp {
    error: Common.Error
    executions: Common.Execution[]
    progress: Common.ExecutionProgress[]
}
interface GetExecutionReq {
   idOrDate: string
}
interface GetExecutionResp {
    error: Common.Error
    execution: Common.Execution
}
interface UpdateExecutionReq {
    execution: Common.Execution
    mode: string
}
interface UpdateExecutionResp {
    error: Common.Error
    execution: Common.Execution
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
    error: Common.Error
    execution_stats: ExecutionStats[]
}
interface GetExecutionStatsReq {
    id: number
}
interface GetExecutionStatsResp {
    error: Common.Error
    execution_stats: ExecutionStats
}
interface UpdateExecutionStatsReq {
    execution_stats: ExecutionStats
}
interface UpdateExecutionStatsResp {
    error: Common.Error
    execution_stats: ExecutionStats
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
    address: Common.Address
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
    error: Common.Error
    deliveries: Delivery[]
}
interface UpdateDeliveriesReq {
    deliveries: Delivery[]
}
interface UpdateDeliveriesResp {
    error: Common.Error
}
interface UpdateDripReq {
    emails: string[]
    hours: number
}
