package entities

type ChangeDataCaptureMessage struct {
	ID        int    `json:"id"`
	ColA      int    `json:"col_a"`
	ColB      string `json:"col_b"`
	CreatedAt string `json:"created_at"`
}

type ChangeDataCaptureSource struct {
	Version int    `json:"version"`
	App     string `json:"app"`
}

type ChangeDataCaptureEventPayload struct {
	Op     string                    `json:"op"`
	Before *ChangeDataCaptureMessage `json:"before"`
	After  *ChangeDataCaptureMessage `json:"after"`
	Source ChangeDataCaptureSource   `json:"source"`
}
