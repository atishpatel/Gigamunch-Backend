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
