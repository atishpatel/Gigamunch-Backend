interface ErrorOnlyResp {
    error: Common.Error
}
interface TokenOnlyReq {
    token: string
}
interface TokenOnlyResp {
    error: Common.Error
    token: string
}
interface SubmitCheckoutReq {
    email: string
    first_name: string
    last_name: string
    phone_number: string
    address: Common.Address
    delivery_notes: string
    reference: string
    payment_method_nonce: string
    servings: string
    vegetarian_servings: string
    first_delivery_date: string
    campaigns: Common.Campaign[]
    reference_email: string
}
interface UpdatePaymentReq {
    email: string
    payment_method_nonce: string
}
