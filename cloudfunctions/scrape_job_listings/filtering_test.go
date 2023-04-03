package cloudfunction

import "testing"

func TestSkipThread(t *testing.T) {
	tests := []struct {
		link string
		want bool
	}{
		{
			link: "https://www.ycombinator.com/companies/aviator/jobs/NCybOOq-growth-marketing-manager",
			want: true,
		},
		{
			link: `https://account.ycombinator.com/authenticate?continue=https%3A%2F%2Fwww.workatastartup.com%2Fapplication%3Fsignup_job_id%3D59974&defaults%5BsignUpActive%5D=true&defaults%5Bwaas_company%5D=24191""https://account.ycombinator.com/authenticate?continue=https%3A%2F%2Fwww.workatastartup.com%2Fapplication%3Fsignup_job_id%3D59974&defaults%5BsignUpActive%5D=true&defaults%5Bwaas_company%5D=24191`,
			want: false,
		},
		{
			link: "https://jobs.lever.co/memfault/730541eb-637f-4d9d-9526-8949432f9a34123",
			want: false,
		},
		{
			link: "https://jobs.lever.co/foo/730541eb-637f-4d9d-9526-8949432f9A34",
			want: true,
		},
		{
			link: "https://jobs.lever.co/foo/730541eb-637f-4d9d-9526-8949432f9A34/apply",
			want: false,
		},
		{
			link: "https://boards.greenhouse.io/supabase/jobs",
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.link, func(t *testing.T) {
			got := CheckLinkPatterns(tc.link)
			if got != tc.want {
				t.Fail()
			}
		})
	}
}
