local charset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
math.randomseed(os.clock())

function randomWords(length)
	local ret = {}
	local r
  local words = {
    'satu',
    'dua',
    'tiga',
    'empat',
    'lima',
    'enam',
    'tujuh',
    'delapan',
    'sembilan',
    'sepuluh',
  }
	for i = 1, length do
		r = math.random(1, 10)
		table.insert(ret, words[r])
	end
	return table.concat(ret, " ")
end

local title=randomWords(3)
local content=randomWords(10)
wrk.method = "POST"
wrk.body   = string.format("title=%s&content=%s", title, content)
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

--login
--wrk.method = "POST"
--wrk.body   = string.format("username=test&password=test"
--wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

--logout
--wrk.headers["Cookie"] = "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyNTQ1IiwibmFtZSI6InRlc3QiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJ1c2VybmFtZSI6InRlc3QiLCJQYXNzd29yZEhhc2giOiIiLCJjcmVhdGVkX2F0IjoiMjAyNC0wMi0wNVQxNDozODozOSswNzowMCIsInVwZGF0ZWRfYXQiOiIyMDI0LTAyLTA1VDE0OjM4OjM5KzA3OjAwIiwiZXhwIjoxNzA3Mzc3OTYzfQ.V4VLsFP0gK2dxJzoxZwWCUY5gtssziy_ct1iufDvck8"
