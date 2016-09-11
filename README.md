# About
This service will implement to invoke docker registry API with flex_auth_service.

# Manual token-based workflow to list repositories

You can skip this section if you're not interested in how a token can be requested 
manually to list the repositories inside a registry.


	# This is the operation we want to perform on the registry
	registryURL=https://127.0.0.1:5000/v2/_catalog

	# Save the response headers of our first request to the registry to get the Www-Authenticate header
	respHeader=$(tempfile);
	curl -k --dump-header $respHeader $registryURL

	# Extract the realm, the service, and the scope from the Www-Authenticate header
	wwwAuth=$(cat $respHeader | grep "Www-Authenticate")
	realm=$(echo $wwwAuth | grep -o '\(realm\)="[^"]*"' | cut -d '"' -f 2)
	service=$(echo $wwwAuth | grep -o '\(service\)="[^"]*"' | cut -d '"' -f 2)
	scope=$(echo $wwwAuth | grep -o '\(scope\)="[^"]*"' | cut -d '"' -f 2)

	# Build the URL to query the auth server
	authURL="$realm?service=$service&scope=$scope"

	# Query the auth server to get a token
	token=$(curl -ks -H "Authorization: Basic $(echo -n "mozart:password" | base64)" "$authURL")

	# Get the bare token from the JSON string: {"token": "...."}
	token=$(echo $token | jq .token | tr -d '"')

	# Query the registry again, but this time with a bearer token
	curl -vk -H "Authorization: Bearer $token" $registryURL
As a result you should get a list of repositories in your registry. If you have pushed only the busybox image from above to your registry you should see an HTTP body like this:

	{"repositories":["anyuser/busybox"]}
	

# Refer to
	https://github.com/kwk/docker-registry-setup