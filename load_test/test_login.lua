local charset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
math.randomseed(os.clock())
function randomString(length)
	local ret = {}
	local r
	for i = 1, length do
		r = math.random(1, #charset)
		table.insert(ret, charset:sub(r, r))
	end
	return table.concat(ret)
end

local username=randomString(20)
wrk.method = "POST"
wrk.body   = string.format("name=%s&email=%s@example.com&username=%s&password=test", username, username, username)
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

--login
--wrk.method = "POST"
--wrk.body   = string.format("username=test&password=test"
--wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

--logout
--wrk.headers["Cookie"] = "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyNTQ1IiwibmFtZSI6InRlc3QiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJ1c2VybmFtZSI6InRlc3QiLCJQYXNzd29yZEhhc2giOiIiLCJjcmVhdGVkX2F0IjoiMjAyNC0wMi0wNVQxNDozODozOSswNzowMCIsInVwZGF0ZWRfYXQiOiIyMDI0LTAyLTA1VDE0OjM4OjM5KzA3OjAwIiwiZXhwIjoxNzA3Mzc3OTYzfQ.V4VLsFP0gK2dxJzoxZwWCUY5gtssziy_ct1iufDvck8"
