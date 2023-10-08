package cavage

import (
	"testing"
	"time"
)

func TestStringification(t *testing.T) {
	params := ParamsWithSignature{
		Params: Params{
			Created: time.Unix(1699239353674, 0),
		},
		Signature: "c/mpkjH6lgJFjBcRdvtEVHjHYoJJK8K8XFdlpbkQU8xIWtqwKd1o6LP8GpaWjRyuYcxCJdWj5lMBsOjTCWCOAW/jMMt/v1l0DIv6O7xHIy3xH9RTHnopsG52hsMvRFbE23j4keaehvZkdxNKmB5o7u7layTeiQ9KAGa7ENUCe2Fqupg3Vqiu9bVEGdT/DTcue9M9kP34wDjduvjK9H9e5gr7RWuNOxihkvyoOHOUpDUlvSDvVFOyluZpCylFXxLHBoXZqP7274eERqrpUvCY5sspEIcRbFDLVK+k3J7+k3w94aVWdPiPQcNTOkuJcpUFrdSRdLKzN61dBbl/4awYVg==",
	}

	expected := "signature=\"c/mpkjH6lgJFjBcRdvtEVHjHYoJJK8K8XFdlpbkQU8xIWtqwKd1o6LP8GpaWjRyuYcxCJdWj5lMBsOjTCWCOAW/jMMt/v1l0DIv6O7xHIy3xH9RTHnopsG52hsMvRFbE23j4keaehvZkdxNKmB5o7u7layTeiQ9KAGa7ENUCe2Fqupg3Vqiu9bVEGdT/DTcue9M9kP34wDjduvjK9H9e5gr7RWuNOxihkvyoOHOUpDUlvSDvVFOyluZpCylFXxLHBoXZqP7274eERqrpUvCY5sspEIcRbFDLVK+k3J7+k3w94aVWdPiPQcNTOkuJcpUFrdSRdLKzN61dBbl/4awYVg==\", created=1699239353674"

	if params.String() != expected {
		t.Errorf("expected %s, got %s", expected, params.String())
	}
}
