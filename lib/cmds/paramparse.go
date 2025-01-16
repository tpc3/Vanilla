package cmds

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func ParseParam(param string, shortParams map[string]any, longParams map[string]any) ([]string, error) {
	var processing, next *string
	var nextQuoted bool
	var unname []string
	for {
		for {
			if processing == nil {
				if next != nil && nextQuoted {
					if unname == nil {
						unname = append(unname, *next)
					} else {
						return nil, errors.New("duplicate unnamed argument")
					}
				} else {
					processing = next
				}
				next = nil
			}
			if next != nil {
				break
			}
			if param == "" {
				break
			}
			if strings.HasPrefix(param, "\"") {
				split := strings.SplitN(param, "\"", 3)
				if len(split) != 3 {
					return nil, errors.New("quote syntax error")
				}
				next = &split[1]
				nextQuoted = true
				param = split[2]
			} else {
				split := strings.SplitN(param, " ", 2)
				next = &split[0]
				nextQuoted = false
				if len(split) == 2 {
					param = split[1]
				} else {
					param = ""
				}
			}
		}
		if processing == nil {
			break
		}
		if strings.HasPrefix(*processing, "--") {
			split := strings.SplitN(strings.TrimPrefix(*processing, "--"), "=", 2)
			switch p := (longParams[split[0]]).(type) {
			case **int:
				if len(split) != 2 {
					return nil, errors.New("missing value: " + split[0])
				}
				conv, err := strconv.Atoi(split[1])
				if err != nil {
					return nil, err
				}
				*p = &conv
			case **int64:
				if len(split) != 2 {
					return nil, errors.New("missing value: " + split[0])
				}
				conv, err := strconv.ParseInt(split[1], 10, 64)
				if err != nil {
					return nil, err
				}
				*p = &conv
			case **string:
				if len(split) != 2 {
					return nil, errors.New("missing value: " + split[0])
				}
				conv := split[1]
				*p = &conv
			case **struct{}:
				if len(split) != 1 {
					return nil, errors.New("unused value: " + *processing)
				}
				conv := struct{}{}
				*p = &conv
			default:
				log.Fatal("Unsupported parse type: longParams: ", split[0])
			}
		} else if strings.HasPrefix(*processing, "-") {
			chars := strings.Split(strings.TrimPrefix(*processing, "-"), "")
			for i, v := range chars {
				if shortParams[v] == nil {
					return nil, errors.New("not found short parameter: " + v)
				}
				if i+1 == len(chars) {
					switch p := (shortParams[v]).(type) {
					case **int:
						if next == nil {
							return nil, errors.New("missing value: " + v)
						}
						conv, err := strconv.Atoi(*next)
						if err != nil {
							return nil, err
						}
						*p = &conv
						next = nil
					case **int64:
						if next == nil {
							return nil, errors.New("missing value: " + v)
						}
						conv, err := strconv.ParseInt(*next, 10, 64)
						if err != nil {
							return nil, err
						}
						*p = &conv
						next = nil
					case **string:
						if next == nil {
							return nil, errors.New("missing value: " + v)
						}
						conv := *next
						*p = &conv
						next = nil
					case **struct{}:
						conv := struct{}{}
						*p = &conv
					default:
						log.Fatal("Unsupported parse type: shortParams: ", v)
					}
				} else {
					p, ok := (shortParams[v]).(**struct{})
					if !ok {
						return nil, errors.New("invalid usage: " + v)
					}
					conv := struct{}{}
					*p = &conv
				}
			}
		} else {
			if unname == nil {
				unname = append(unname, *processing)
			} else {
				return nil, errors.New("duplicate unnamed argument")
			}
		}
		processing = nil
	}
	return unname, nil
}
