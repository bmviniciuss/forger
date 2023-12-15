package core

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type InterpolationFunction struct {
	regex  *regexp.Regexp
	prefix string
}

var (
	requestVarFunction       = InterpolationFunction{regexp.MustCompile(`\$requestVar\('([^']+)'\)`), `$requestVar('`}
	requestParameterFunction = InterpolationFunction{regexp.MustCompile(`\$requestParameter\('([^']+)'\)`), `$requestParameter('`}
	requestHeaderFunction    = InterpolationFunction{regexp.MustCompile(`\$requestHeader\('([^']+)'\)`), `$requestHeader('`}
)

type GeneratorFunction struct {
	regex     *regexp.Regexp
	generator func() string
}

const (
	utcTimeLayout = "2006-01-02T15:04:05.000Z"
)

var (
	uuidGeneratorFunction = GeneratorFunction{
		regex: regexp.MustCompile(`\$uuid`),
		generator: func() string {
			return uuid.NewString()
		},
	}
	iso8601GeneratorFunction = GeneratorFunction{
		regex: regexp.MustCompile(`\$iso8601`),
		generator: func() string {
			return time.Now().Format(utcTimeLayout)
		},
	}
)

type Interpolator func(input string, r *http.Request) string

func makeRequestVarInterpolator(input string, r *http.Request) string {
	return interpolateWithFunction(input, requestVarFunction, func(paramID string) string {
		return chi.URLParam(r, paramID)
	})
}

func makeRequestParameterInterpolator(input string, r *http.Request) string {
	return interpolateWithFunction(input, requestParameterFunction, func(paramID string) string {
		return r.URL.Query().Get(paramID)
	})
}

func makeRequestHeaderInterpolator(input string, r *http.Request) string {
	return interpolateWithFunction(input, requestHeaderFunction, func(paramID string) string {
		return r.Header.Get(paramID)
	})
}

func InterpolateString(input string, r *http.Request) string {
	funcs := []Interpolator{
		makeRequestVarInterpolator,
		makeRequestParameterInterpolator,
		makeRequestHeaderInterpolator,
	}
	result := input
	for _, f := range funcs {
		result = f(result, r)
	}

	generators := []GeneratorFunction{uuidGeneratorFunction, iso8601GeneratorFunction}
	for _, g := range generators {
		result = processGeneratorFunction(result, g)
	}
	return result
}

func interpolateWithFunction(input string, function InterpolationFunction, valExtractor func(paramID string) string) string {
	matches := function.regex.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 || matches == nil {
		return input
	}

	funcParamsMap := make(map[string]string)
	for _, match := range matches {
		if len(match) == 2 {
			parameterName := match[1]
			parameterReqVal := valExtractor(parameterName)
			funcParamsMap[parameterName] = parameterReqVal
		}
	}

	result := function.regex.ReplaceAllStringFunc(input, func(match string) string {
		paramName := strings.TrimSuffix(strings.TrimPrefix(match, function.prefix), "')")
		if val, ok := funcParamsMap[paramName]; ok {
			return val
		}
		return match
	})
	return result
}

func processGeneratorFunction(input string, generator GeneratorFunction) string {
	return generator.regex.ReplaceAllStringFunc(input, func(_ string) string {
		return generator.generator()
	})
}
