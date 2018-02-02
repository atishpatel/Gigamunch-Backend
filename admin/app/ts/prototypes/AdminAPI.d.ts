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
