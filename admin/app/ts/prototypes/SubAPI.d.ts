declare namespace SubAPI {
interface ErrorOnlyResp {
    error: Common.Error
}
interface DateReq {
    date: string
}
interface GetUserSummaryReq {

}
interface GetUserSummaryResp {
  error: Common.Error
  // is an active subscriber
  is_active: boolean
  on_probation: boolean
  // has finished checkout flow
  has_subscribed: boolean
  is_logged_in: boolean
}
interface GetAccountInfoReq {

}
interface GetAccountInfoResp {
  error: Common.Error
  subscriber: Common.Subscriber
  payment_info: PaymentInfo
}
interface PaymentInfo {
  card_number_preview: string
  card_type: string
}
interface GetExecutionsReq {
  start: number
  limit: number
}
interface GetExecutionsDateReq {
    date: string
}
interface GetExecutionsResp {
  error: Common.Error
  execution_and_activity: Common.ExecutionAndActivity[]
  activities: Common.Activity[]
}
interface GetExecutionReq {
    idOrDate: string
}
interface GetExecutionResp {
  error: Common.Error
  execution_and_activity: Common.ExecutionAndActivity
}
interface ChangeActivityServingsReq {
  id: string
  servings_non_veg: number
  servings_veg: number
  date: string
}
interface ChangeSubscriberServingsReq {
  id: string
  servings_non_veg: number
  servings_veg: number
}
interface UpdateSubscriberReq {
  first_name: string
  last_name: string
  address: Common.Address
  delivery_notes: string
  phone_number: string
}
interface UpdatePaymentReq {
  payment_method_nonce: string
}
interface ChangePlanDayReq {
  new_plan_day: string
}
}