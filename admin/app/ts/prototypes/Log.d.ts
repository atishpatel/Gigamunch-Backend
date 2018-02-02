interface Log {
    id: number
    log_name: string
    timestamp: google.protobuf.Timestamp
    type: string
    labels: string[]
    severity: number
    payload: string
}
