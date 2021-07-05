package gotau

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTAUTags_UnmarshalJSON(t *testing.T) {
	type tmp struct {
		TagIDs TAUTags `json:"tag_ids"`
	}
	tagString := "{\"tag_ids\": \"['6ea6bca4-4712-4ab9-a906-e3336a9d8039']\"}"

	data := new(tmp)
	err := json.Unmarshal([]byte(tagString), data)
	require.NoError(t, err)
	require.Len(t, data.TagIDs, 1)
	require.Equal(t, "6ea6bca4-4712-4ab9-a906-e3336a9d8039", data.TagIDs[0])

	tagString = "{\"tag_ids\": \"['6ea6bca4-4712-4ab9-a906-e3336a9d8039', '621fb5bf-5498-4d8f-b4ac-db4d40d401bf']\"}"

	data = new(tmp)
	err = json.Unmarshal([]byte(tagString), data)
	require.NoError(t, err)
	require.Len(t, data.TagIDs, 2)
	require.Equal(t, "6ea6bca4-4712-4ab9-a906-e3336a9d8039", data.TagIDs[0])
	require.Equal(t, "621fb5bf-5498-4d8f-b4ac-db4d40d401bf", data.TagIDs[1])
}
