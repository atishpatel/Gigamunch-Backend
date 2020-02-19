declare namespace SubAPI {
interface ErrorOnlyResp {
    error: Common.Error
}
interface DateReq {
    date: string
}
interface GetUserSummaryReq {
undefined: undefined
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
undefined: undefined
}
interface GetAccountInfoResp {
  error: Common.Error
  email_prefs: Common.EmailPref[]
  phone_prefs: Common.PhonePref[]
  address: Common.Address
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
}