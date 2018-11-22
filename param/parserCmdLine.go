package param

import (
	"errors"
	"fmt"
	"github.com/NickWells/golem/location"
	"github.com/NickWells/golem/strdist"
	"strings"
)

// Remainder returns any arguments that come after the terminal parameter
func (ps *ParamSet) Remainder() []string { return ps.remainingParams }

// findClosestMatch finds the parameter with the name which is the shortest
// distance from the passed value and returns it together with the distance
func (ps *ParamSet) findClosestMatch(badParam string) (string, int) {
	var minDist = -1
	closestMatches := make([]string, 0)
	var closestMatch string

	for p := range ps.nameToParam {
		if minDist == -1 {
			minDist = strdist.Levenshtein(badParam, p)
			closestMatches = append(closestMatches, p)
		} else {
			dist := strdist.Levenshtein(badParam, p)
			if dist < minDist {
				minDist = dist
				closestMatches[0] = p
				closestMatches = closestMatches[:1]
			} else if dist == minDist {
				closestMatches = append(closestMatches, p)
			}
		}
	}

	sep := ""
	for _, match := range closestMatches {
		closestMatch += sep + match
		sep = " or "
	}
	return closestMatch, minDist
}

func (ps *ParamSet) reportMissingParams(missingCount int) {
	var err error

	byPosMiniHelp := "The first"
	if len(ps.byPos) == 1 {
		byPosMiniHelp += " parameter should be: <" + ps.byPos[0].name + ">"
	} else {
		byPosMiniHelp +=
			fmt.Sprintf(" %d parameters should be: ", len(ps.byPos))
		sep := "<"
		for _, bp := range ps.byPos {
			byPosMiniHelp += sep + bp.name
			sep = ">, <"
		}
		byPosMiniHelp += ">"
	}

	if missingCount == 1 {
		err = errors.New("A parameter is missing," +
			" one more positional parameter is needed. " +
			byPosMiniHelp)
	} else {
		err = fmt.Errorf(
			"Some parameters are missing,"+
				" %d more positional parameters are needed. %s",
			missingCount, byPosMiniHelp)
	}
	ps.errors[""] = append(ps.errors[""], err)
}

func (ps *ParamSet) getParamsFromStringSlice(source string, params []string) {
	loc := location.New(source)

	if len(ps.byPos) > 0 {
		missingCount := len(ps.byPos) - len(params)
		if missingCount > 0 {
			ps.reportMissingParams(missingCount)
			return
		}

		for i, pp := range ps.byPos {
			pStr := params[i]
			loc.Incr()
			loc.SetContent(pStr)

			pp.processParam(source, loc, pStr)

			if pp.isTerminal {
				ps.remainingParams = params[i+1:]
				return
			}
		}
	}

	for i := len(ps.byPos); i < len(params); i++ {
		pStr := params[i]
		loc.Incr()
		loc.SetContent(pStr)

		if pStr == ps.terminalParam {
			ps.remainingParams = params[i+1:]
			return
		}

		paramParts := strings.SplitN(pStr, "=", 2)
		trimmedParam, err := trimParam(paramParts[0])
		if err != nil {
			ps.errors[trimmedParam] = append(ps.errors[trimmedParam],
				loc.Error(err.Error()))
			continue
		}

		if p, ok := ps.nameToParam[trimmedParam]; ok {
			if p.setter.ValueReq() == Mandatory &&
				len(paramParts) == 1 {
				if i < (len(params) - 1) {
					i++
					loc.Incr()
					paramParts = append(paramParts, params[i])
					loc.SetContent(pStr + " " + params[i])
				}
			}
			p.processParam(source, loc, paramParts)
		} else {
			ps.recordUnexpectedParam(trimmedParam, loc)
		}
	}
}

// trimParam trims the parameter of any leading dashes
func trimParam(param string) (string, error) {
	trimmedParam := strings.TrimPrefix(param, "--")
	if trimmedParam != param {
		return trimmedParam, nil
	}
	trimmedParam = strings.TrimPrefix(param, "-")
	if trimmedParam != param {
		return trimmedParam, nil
	}
	return param, fmt.Errorf(
		"'%s' is a parameter but does not start with either '--' or '-'",
		param)
}
