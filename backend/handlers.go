package main

import (
	"github.com/valyala/fasthttp"
)

func GetQueriedDomains(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Hello")
}

func PostDomainAndGetInfo(ctx *fasthttp.RequestCtx) {
	domain, err := ctx.UserValue("domain").(string)
	if !err {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
	md := &MethodCaller{}
	country, organization := md.ObtainServerCountryAndOrganization(domain)
	ctx.WriteString(country + organization)
}
