# TODO List

## High Priority
- [ ] refresh token reuse detection
- [ ] unit tests for membership and family store in sqlite
- [ ] "real", not mocked API tests
- [ ] revise testing on different levels (gaps!)
- [ ] doc testing architecture
- [ ] service methods too long, extract into funcs
- [ ] 

## Medium Priority
  
- [ ] Cookies with refresh token (HttpOnly)
- [ ] Add transactional integration test (SQLite)
- [ ] Add JWKS endpoint for gateway

## LATER Nice-to-haves
- [ ] OAuth / OpenID Connect
- [ ] Docker multi-arch image (in case nodes have diff architecture: amd64, arm64)
- [ ] Multi-factor authentication
- [ ] Social login
- [ ] Token revocation lists
- [ ] Token introspection endpoint
- [ ] Fine-grained permissions
- [ ] refactor register svc unit test, extract func to create the svc