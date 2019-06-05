declare namespace AdminAPI {
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
    id: string
    reason: string
}
interface ReplaceSubscriberEmailReq {
    old_email: string
    new_email: string
}
interface ChangeSubscriberServingsReq {
    id: string
    servings_non_veg: number
    servings_veg: number
}
interface DiscountSubscriberReq {
    user_id: string
    discount_amount: number
    discount_percent: number
}
interface GetSubscriberDiscountsResp {
    error: Common.Error
    discounts: Common.Discount[]
}
interface GetSubscriberActivitiesResp {
    error: Common.Error
    activities: Common.Activity[]
}
interface SkipActivityReq {
    email: string
    date: string
    id: string
}
interface UnskipActivityReq {
    email: string
    date: string
    id: string
}
interface SetupActivitiesReq {
    hours: number
}
interface RefundActivityReq {
    emails: string[]
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
interface ChangeActivityServingsReq {
    id: string
    servings_non_veg: number
    servings_veg: number
    date: string
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
interface GetLogsForUserReq {
    start: number
    limit: number
    id: string
}
interface GetLogsByExecutionReq {
    execution_id: number
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
interface ProcessActivityReq {
    email: string
    date: string
    id: string
}
interface ProcessActivityResp {
    error: Common.Error
}
interface SendCustomerSMSReq {
    emails: string[]
	message: string
}
interface GetHasSubscribedReq {
    start: number
    limit: number
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
interface GetHasSubscribedRespV2 {
    error: Common.Error
    subscribers: Common.Subscriber[]
}
interface GetSubscriberReq {
    email: string
}
interface GetSubscriberResp {
    error: Common.Error
    subscriber: Subscriber
}
interface GetSubscriberRespV2 {
    error: Common.Error
    subscriber: Common.Subscriber
}
interface UserIDReq {
    ID: string
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
}