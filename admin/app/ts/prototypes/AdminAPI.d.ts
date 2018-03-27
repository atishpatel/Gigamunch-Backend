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
