package assets

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"sort"
)

// AssetNames returns a list of all assets
func AssetNames() []string {
	an := make([]string, len(data))
	i := 0
	for k := range data {
		an[i] = k
		i++
	}
	sort.Strings(an)
	return an
}

// Get returns an asset by name
func Get(an string) ([]byte, bool) {
	if d, ok := data[an]; ok {
		return d, true
	}
	return nil, false
}

// MustGet returns an asset by name or explodes
func MustGet(an string) []byte {
	if r, ok := Get(an); ok {
		return r
	}
	panic(errors.New("could not find asset: " + an))
}

func decompress(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	r, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	defer r.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.Bytes()
}

var data = make(map[string][]byte, 82)

func init() {
	data["internal/bootstrap/00_builtins.lisp"] = decompress("H4sIAAAAAAAA/3RVy3LbMAy8+yuUW50Op9Nr1fRfYBKy0VAEDVKu1a/vUFYzE4A+cvFarpbQOI7jABGHE3MtVSD/GE4LxUrpcPgScHL7aQh8/AzQpICIVSFTOqoughBUEt4g6tbF4fV4OIzjcI58glg0GfQRBI8anR5FJyjkhywYyENFXU3FJerMhMqzRd9x/cMSHp1n8MKWzbThqva6cNUUy5oq3F0vtLXAe4YU3PfnMUtwn93olXU+sVGrrPpaZ0wWpOIe9RaP7Hsf6bpApIlwF6dUoXQ206toEaq8dGa3vK0NXhdM3ny20jzx2V8kRXtO0EARDOQ5FQv9No5OCkn1YqTUvaEU9roTGQo39JVFP5AmhLln+wKGChVHKeAdO5boUdicknMvvcOOivsgOI5DWmYU8orsV1Xk1PlVnb9pW3PQ939TKS8a+KXPOuGnPr91NM7c1NMbrC0GPD8JQHpIMS3JV+Kk7ZlBKplXAjnHtfOB+nDJ6LcebZDn5BcRTH5Vo86slfUX0P7IwjMVvWXa3f8H2pAIf9enL64Fmx3NNE4ezLKnWFE/9RmyQiq8m70trLMEw+J13sTiEPzl9UFdgAoOkMIg6PmGov81YK+/Z+4Cb+rWodKMqnYPuRY6Hv4FAAD//yFBUBAmBwAA")
	data["internal/bootstrap/01_basics.lisp"] = decompress("H4sIAAAAAAAA/3yOQWrDMBBF9z7F7ybIlEDbZU3B9yiFTqQxGCyNo3Gy1NmLHFFsg7IT+u89puu6DjQxLiKLLpHmT1xIR6tNYxwPnmwUOB5CA3wH8owTBolefxrgNyNI67cZQnmlfgXatt02OEaJOXKCneimXAqRRmUYUhWL1JftIM8URluX13ZV5ms27zTd8vFW/FzsUc98RXosqV+XvfpSdYMsqARKISDPdw7/hWyaLxgvDg/no8XbjhfnnuPvGzzYPftasA3j+MCcN8xfAAAA//81dvsl/QEAAA==")
	data["internal/bootstrap/02_functions.lisp"] = decompress("H4sIAAAAAAAA/2xSTa+cMAy88ytGe9iGA1V7fXt5/wOhbh6YbaqQ0MQ8dS/89ipfkK7qAyL2eDJj53a73SA14cNa9uzk+oZ5MyMra3zTiInmRY7OQnpPjjvpHr4BRD/AKN3G31HLzdPQABBOKk+4VGiksseyecYHYZXK0XRpU/O11GP/XYzWTOEPwC5K67cWoqbc34Ujz0jfjGojY4g30p6QI0s6yb63IZI3A2M7uzZAPxwpNZFhxc+Q/ZQ6CPuU+iiP1niWhnUNELNBf8WPIUIr7LJ2dmNyARsGi2gg6O7CMdoWaoZQvvP0u6okN+KkgJiV81xDILTyjIO4zSOpIJnmBVckpuWGO/JWi/m0ndAwxLbzjCsO5RCaGH30VOb9IOOfCy4heSkrAaCMIVdAUc4Xua76Wak6wdHpv7arocSIHhMiu06IIUHueSV7YB6Otr0eaH1LkneSVq8kDenX9ofDFiv7lfsX2+nC/2vo93exyDW9mXl4GUbeTzHzIuJrUMDSPYixEP+0E64o7HchQn7P9T0BWuzvibP5GwAA//+yPmM87gMAAA==")
	data["internal/bootstrap/03_branching.lisp"] = decompress("H4sIAAAAAAAA/5RU7YrbMBD876eYphAkaOB6PxsK9x4hEMVWEoEsBUlOuT9+9qJdfzuh14NL4t0Z7WjWu/v9fg9lNc7ep5iCuv/COShX3oy7FoWo9KVWZfBwPhXA4aHssQBOwlzQPpTFRdmokUKjpZziG2d1jAUgDknHdIQzVg6PSDftjt0xFHDGos1RuUBB26hn0DZHevCs6J+bdi9LXnyo5+fkCIEmqC3h4gwoKo/2g+KS8auqO/bni5Xpsjn6r8oZOK2+qKxcRfwjN4B+l1Y12TD+5ti2e4q5dzgJqxMOylXf0YougzdJSSBXp5xQrkL7IUJW0sGkpNxSyDdWMi90Es6n/pSBP+P50OkfDPuyfB9eqs+p/C98WOtfSScNT5UT/bnw0g/W/4d0Vk6dBSbaWXr+49kbcz+He/FAsZ3CxJ1KvkY+SvLneAaAyTs0jwMtl1iGRb7R4NbKszmcj1g6aS47q3kKzsZVxl0nY55TaPs4De9k9KZ4HniyS8WoQ9qpcI0FNzfuHrpMPqBjSGx6bpeom5hw1ojN/W6NrjbM/I33bL8biC+YpXdJGYd3aKtr7VLcyLF18bMGRM97kz9GX3hT9Kmhb93b2mYm9ePYU6hDHM63pq22ejtptyxdHffJwtXJSuljW5x99fkMzXslZ6no3wAAAP//Gdd6/AoGAAA=")
	data["internal/bootstrap/04_concurrency.lisp"] = decompress("H4sIAAAAAAAA/2xPy2rEMAy8+ysGAkXuJySX/kcI1Oso24XE2ir2YS/99rJ+tEmpLx5mJM3MMAwD3Mq4iMQ9qrv38BJ8UuXgH8bQzMvmvAquYoDxBReZH5MB3ukqr6AlYJzw9fakrbWnDQ6sLvLfvZUjRv/hQgfQ87cG7flVdu5AfQbIUweZt1sEQH0GRZ2KTDlggdlBeU9r7ECztHzT7yUqTofbdd5WivqdP6vFudeSYtL/W91VtlsuUNEh3Q93itTsmmat+Q4AAP//b4/Cv5UBAAA=")
	data["internal/bootstrap/05_lazy.lisp"] = decompress("H4sIAAAAAAAA/4RUXW/iMBB851eMWulk9xQJ7rFUVf8HQqpJFhrV2MFr2uMe+O0nf8WQUrVPjWdnd2c8ZrlcLqE0YWOtZ+/U8Ait/p3AdDiSaYlnM9HRdq9aZyPSMB1mwOoXNrY7rWfAqyjnDxBbg9Ua55cASikv2d42itm2ic104MRWw6BPiBBEa02rPM4vAf/aQPfsb/MD8iP9g1pv3e0GCfuuhYFX79R8vvWaAn9w1KG1WocW4sIYiM83Mo0mjxWPFeUcIhLFtnfswVImDHEs1/NynME6G4nvKNHj36XIrXUNqfYt7Mh0aOjv4Bj1soRiJucb5XYct+25oY+wmCaDkSEl7korbHrT9WbH2B/ZY0NQeKeoN5t2J6MJQbJRe7qHqKPnVQvT4R6X2EImb0S/hXieroA/hZpaB833WXo2oK67Hqe8irL4A85xZAzlOW62xojinBqOWS3TfmhQ66fGf+f5q9iRIac8JbV1g1otaN97iM5+7W8wKOd731sTbF619mh8ShbECKEcFzDqKdXsaahhvMprCWbP8UYjdXYVyvzuUgynEy5SOt0lDhWds0OdfxVZA6fMjuKeabH4jTlMr7HICrRif42GkxKa+P9cYoFmITMjvaLKu6iMUN4796vVWFwjlXGzeVSVDQ0BbfdDtKyLZgYJoUomG70yXp/g3ZEmrzu795RMmsvw9Xyr5JE00/j1VDIf5YXZdbPJz0oCxlsv8sTvjITBskqqTwGrdbyt/wEAAP//Xh0tJScGAAA=")
	data["internal/bootstrap/06_threading.lisp"] = decompress("H4sIAAAAAAAA/9xUS4obMRDd9ykeMxCkhSDJMgbje5iG0birbRG11Khkh9no7EEf/8ZODFmZeGPxqrreR6IWi8UC2hLevY8cg55/IO4C6cG4bdeJgUbXABW9soZjB6yjDluKfQeIvbPEDGG4VFFLssPpJ/6AN0RWmklvgoda5pnrg7Z76lH+5BnBF4w+TJyJISxFrC3EtTyI0QSOtVFe8I3Hir0AA0SggvUVfBNqCZFGpMqYVkEirWpXmyk/a34G0U31KjTlj1Szn+gJ0i5jD9q+Ntn9qVv82pFTzsfytpyxWeDr1cOqHorxOuL+Zd0z/l84P995aXhoXfPljcPp6Z75DP9DAtfGUply4ywrQLJo5Vu9x50z6Z+kNt4NamP1nilvHf6YysoZHdYVrZMr4xxogHBxh1rD18vEMsVV9dtRMN5EmXjQtocwI1KdlPhjyuEg5W9lPt4mWhQ+ekuVsWWpmSlEpcOWq3jDig7ksgt3bJUSL+2Iac8R74RZm0DDSzFVtlQN9+QxrcSkZ4jPyUEtJcSsQzTReIfvZ5K7Zp7ezV/t/A4AAP//8rO3edEGAAA=")
	data["internal/bootstrap/07_exceptions.lisp"] = decompress("H4sIAAAAAAAA/4xW247jNgx991ewWaArPQTt9nGCAvMfQYAqNj3jriJ5KGV23aL59oK62JKd7E4e5iJRPIeHR1QOh8MBlEY4W+udJzU+AX5vcfSDNa5pRIe9gcHtW6V1A3B00wVara4OTw2AUKYDMbi9tq3S4KaLbCB+wurgfIpe1vGN40D0Ay27UsoKzLev+/NgusG8MGpv6VLhvWPrLQGvL5n/hD9AaDRxuWYS+AnegC/30BglcoFRkcOuguPq4XMI3NSz5gvC+BzFWHOgsT5EO3wD8XTWtv2asOSKUT8YpfX0EU4pdMOqBstRM9yDuFjhz6J+RB2/j3Sf9y+WPqblnbIWCE/TfiQc0XQME/PDV5y+WepKZ7bW/J3g4Zj3edXl3ogX9FAnkPK0guLtpR4XUl+NRudmQdLWXAUA/Js0EhKSpvxXLkrI/4pYodHDsa8vRJUtfggE4Y8CRhAzZSB5qgK48m51ZLYu9DBKXijUhTFT79dQYrFoPLk+N2/eO8kOyYCbk1G2zbkn1A6rRKQGh7C7KM13GrvQrXgPE/qusuayPRJ2Q6t8aCv/A0i0d1McMKEZGoR/JVTd3ts4xDhuITU3Sy9ruT06CR+t1qffCQJo5tRqRVjQmtu6IXwmZcoB5SrCyjkkv1f04pqkcGlK2MUmxiTAF3lUA2G3k6Fc9c/E8c1sxHw5HrnxXVEoT8wTDn4vtkN3IWlRGlZmQwZBmux9Fvez0vibRl81/chARaHFTq/YDVHYcLazAVeeMhGx0XV1eXLmBx5JYY80r1T79ormwSz4mKIQ3LXS9Mtyf0vBqtJmJ0NpY3kvNLV/VYy8X/563oUqkIgnpmEX75Bolwj+JXoDxxsSVYSLWXN7Fmoc9QSh2XdaU7CShS7x0h89XRFC/jXZzTMZ6V7UGBrDvNKDkMWvDRmwcj9S5otqyXJ2zvorlHl/etM2IwgI364DoQPlQaNyHqzBdCJdQNY2PULl/N7YJE5GWD+/834AZp/Xj/gys/Ljs/4ykNsYmBC2n1il1r4jRQlP84XrLNyeA7o8SbgtnczgsnI1EsVcn6oJAUDo8kZh89uzKJqa+RbfDfqcMYx+TiLjT+7k/wEAAP//jkgyj8kKAAA=")
	data["internal/bootstrap/08_io.lisp"] = decompress("H4sIAAAAAAAA/8xTTY/bIBS8+1dMt9LqEcntdo+Jtt3/YUUqwc+7qBgcIIlyyW+vgMQfaXvYW31AAubNmxmeN5vNBtIwds7FEL0c1tBfXVVRy53F4OteDvVJx/faalMBTXewCoH32wqgXg4VAOosmqM0W5DuQDokMI7SiLSAck3aCpHggfdCTC3QPKJzvg+Z0nBEE3gPumuOEP2nAhQJibFZQifOfAjQF6zcIa6wPnkdGdRpH2LpWjDUOV+zVO9o2HAP8nwFbK8ct++eaxUGqXiFVCbE3IVd2pDDYM7JXBFc/YXK8sloy6s5i7bxA3H8v2loG80/EkkWPxSK1J7b+sgqOv8jjWCatMxo2+y4XJV5oxdQ79oUnC0nzwJPM205PzewrZVxged8OWsFWuebXL2MtlhQAgrW1W4YaXupvJuo0ey0bbV9C3jEzrXnIjcE9rGW/i3kn2ZpDLcSgYeJaOTpDyFix5D4xef6KM2BUSofco7K2bZofSneJ7qn8sw/qXW4vCY513f//gf0+QbNUVxobP8kMNt9E/PByGl9HmfkLmAsWKaRoujPEwvNsru8lgks6yhObBfqS1mnrTTmDCoiRP6q3wEAAP//+rjaWdYEAAA=")
	data["internal/bootstrap/09_os.lisp"] = decompress("H4sIAAAAAAAA/3SOvWrFMAxGdz/FhwNFHkLTte7Q9wiBurEDAf8EWVn77CXXuT9wuVo0nHOQrLUWLgb8liJV2G2fKFUp8mFJbuYCWVNQwPiGpXCqkwJ+KAbBWMWxdADNO3PI0h+qUbgOh7pH6UC+4O/7Uj/QkH137Fe137nxvqntmpmaQesC+mrSxzAMt5A2XrPE3JDOVT8jer+HRwudqjand35tjPoPAAD//05AX3gcAQAA")
	data["internal/bootstrap/10_predicates.lisp"] = decompress("H4sIAAAAAAAA/6SSwa7bLBCF93mK888vNXjhF2gWfZAo0uXaQ4SKwQFym3SRZ68AJ75VJpWqbhhgDof5Bna73Q7aMd5DyClHPX/FHHm0g86cNhs1svHrTu/0z2uf+LQB9ubsBwzBucMGUJ8yUNZA2VSWVdCVzaLhjL2BqieVsTHllu8OTdGOmvsCUEPwCQbquQQ0l8gPk249t1VGu8SPva3K8VyWXWOa9BADRjb9ajyH9MDyeuKGVUs++6EvW1DpOkGlHKsC9I3utb+1Vt1W7f62bYhfcNuWMldIp1N+FCux3Rpcpb/bNJOu+zOF5+NfUNB/9I8kxdmb/6GMx/5ygPIhQ7X6L92nl12wRdxq8BL2Ne0T6Zsaw/LZnh53aeqtaDtR5Pn4m+h+8yqBTb3OYQKVkToh7a0DeevEZO15ylFO8gmU+CQmHXuQYy8mrR/5wiNomYginVIYQDWIgknPczFpUa7Cpgwqo5j+4CGHCGqRxPbNMUw2MWiZvFCFQmUKkZHbzMcm6F8qtAd57eUL2m+iGuR2zbO7gmqQn2vmwWoHWibyPd/5+iPEEbRMZKvr9B6KU41y58NQ7qpBFJzO2lljywM+pnJJ/FH+UhlFozCOoDCWw78CAAD//z9A478fBgAA")
	data["internal/docstring/and.md"] = decompress("H4sIAAAAAAAA/3SPMW/CMBCFd/+KJ1iApkFtt24MzB07Rqfkgi05dnS+EPrvK8dUoQOZ4tO79323xY5Chz7KcDjsMbLk3wRCslH0tXXSTk5duOBIoTua85X8RMoJahkl3Esc4LlXaIS4i9Uap4QUYwAlxMDgdS2C0JNP/IM85Aqz8x7COkmAWtIyr/GllmV2ieG0hEaJLXOXS+6N2SybBL7polMbs93iFHC+0TB6Ngblxt0L3vCOj71B+RaJv8cmRL07cLfZG/P9INUs0aZC4CvLI7r5v9dUmC0LU/p8zlWZVuwzZL4pqRTG2l+b3wAAAP//MxDI/rIBAAA=")
	data["internal/docstring/apply.md"] = decompress("H4sIAAAAAAAA/1TPTUrAMBAF4P2c4kEXJggFFdy76A3cl7GZ2kCaxPzUentJgwWXM3zvDTNAcYzuB2v1C7J8abTZSganz7qLLxklgC9QbPA0HewqF8komyCmcFgjpmWr+EXA3twdtmQ0LlfJP//XBxVS9wt/ONHgHrvPj0TDgDeP6eQ9OiEClJEVJx7UE57xonVb9T8ecWqi981mSPf4ts4hSanJY36dR/oNAAD//8CDF8f1AAAA")
	data["internal/docstring/assoc.md"] = decompress("H4sIAAAAAAAA/2zQQWvjQAwF4Pv8iodz2YRNWNjbshR6KP0DhZ6VsSYWGWuMJKfxvy91Am2hxwcf7yFt8IvcW8b/My+4UJ35YbfbIhtTsIOg/IaVCIVcGB4255iN06vUeoc/uCM592iKGBiTtYv03OPMy35dwURi/hvNYByzKUjB4xTLtxYp0IbSbHSQfRYd8DKIo8yaQ5pC/MP1UgobayAGui1/LasSbFThiwZdwdfM00oDEsikODJiPaiHKAhFzGOfK7mjkA/S9JDSZoNHxdOVxqlySrj/MAHAP6WR0VHl7pbpxMCfw99bqqQndM+t26b3AAAA//8VhAXUewEAAA==")
	data["internal/docstring/chan.md"] = decompress("H4sIAAAAAAAA/2STwW7bPBCE73qKQXKRjMQP4FvwI6f/VvRWFNCaWllEKdJZLu24T18sZTtW65OhHe7OfEs+o3UTxQ5OmJQzCCX6wQs79SlSgJUjh+bt9g/eVAMpIasUp0UYOpFaoWQeoAkHjiykDEKg3xdk/igcHSONOFEonLf4PjGEcwm6tJwoT68zHeFSzD6rjweTU0TPs9ceY4nV1QsIvQsp8+pbHED3QVt8YxXPJ+tCERx45qgYJc3Qib8czXTBZh+S+7V5wZl8nTsmqarIn7oYtlR7hjlRHpDkrrlxWQTV17DFuwnr8K/zdNeefQigkBPqZJSoPtRu+zKOLEZkDCVPPIAMzhVUGuHVyFY4J5ZVlm3TPD/jv+uM//mSm2azy/yxgf1WsRe3tz0NzWZn0arQcNWYcsdrg+34mGRG31p5idX1zWZXM9tR+jqgaUGxQtS39VvXV6NvEe+fNB8DNw3QBlb8cNP1Rv5szEp7SGi3cBOqPTyNKT11tVTLj6U9yb+lZV7XLN00vZ7YaRK0hgVu6qx2JZYhfBTOdk/MdUjn18AnDihxYDn6GH085BuL/gbv8RaeJ+8m256PmFNWOMqcq/4oPLIIDzjTZfVIVk/k+jJ4IQPapxPDpRIGW5jwWWw1cVeZPQS6dVuS/k2rfWDUdc2fAAAA//8pTMTx+AMAAA==")
	data["internal/docstring/concat.md"] = decompress("H4sIAAAAAAAA/0yOzUrEQBCE7/0UBTmYIAr+3kV8A28iy5AtzUBvzzrdcV2fXjKB4LW+r6q7Qz8WG1PA+XU5QNNv1jPWjJaCvpCZNtLlubIlafHOG8BpKs6lFLRAdsREVPqsgfLxby3bZ2NUHmjhC2UaJxxr+c577rfNa5Guw5Ph5ScdjkoRoI9ypdlj+/ntBre4e8dFf48HPA7DIPI6Zccpq6Iy5mrtXmvt+qZjdXfyFwAA///8jynh/AAAAA==")
	data["internal/docstring/cond.md"] = decompress("H4sIAAAAAAAA/7SSMW/bMBCFd/6KB3uxAsN2hgwJghQdWqBzuytn6mwSoUmBd6rkf1+clLhFgY71dsZ7j/fd0xobX3KH575yBw2cX+7uwEn4U4Oe66nUi8AkUWPJlHCslH2I+ey+lgomH7D/w72HTzQIb22E/R89KWOMKeHI4J+UBlLutqDcIZ4QFVGgddBwxSYXRXuiJNxuMQ85praZ0/b2wh62k1luUXNSZR1qttyigesYhWdT5knfdzJTX4tnEe52zq3X+JzxZaJLn9g5YNPxCRMeH5t5MmwHAJtnTHg4NMBqJEFiEWigjIfDahG8YML94dAsgnNlUq6L5v7wLvr9m0Ux48g6MudV49y3DA1R4Mlu1/6taG8H/ODc4XuK56DpikupDGLRwBo9pXRFn5gk5jPGMqTOfFpgJ6CMJ6sXb3wdS+2e/gfp8sI/SH8EgyApRky6YI+lvon1c2Q/V2XVvc45rx+r3goXg1k+mJ37FQAA//9CQizNwgIAAA==")
	data["internal/docstring/conj.md"] = decompress("H4sIAAAAAAAA/2SQvW7jMBCEez7FAG7OuIOAy3/rIm8QIEVgWCtyZTGmdhVyZcV5+kCR7SYt55vhzqzwx6u8o/AHWs393zUohAJO3LNYgSloVkcWz27zW5vtGoWaxFeuwksXCxru6Bg1Y4opIcS25YzAA0uIsocKrOM5wGc2hp0GrrBBisUWy5B/4H8gHNnbJYmG5XXqYmKQgEpRH8nikdHTgQtEsR8pkxhzATU6GjQHzlH25+vaUbxFlSVT1DBpPmCK1iHR1+napqCMvgMVqHDBkDWMngOaE+qehhqaUbcxGee6cm61wkbw/En9kNg5nBd++48b3OJui3s84BFPa+de558z25iXKc4ldxf2gm53lfsOAAD//0Od9nWqAQAA")
	data["internal/docstring/cons.md"] = decompress("H4sIAAAAAAAA/3ySvY7cMAyEez/FAFvEFg5r3OWnT5EuZYCUB57EtQnYlEPK+/P2gezsFkFypYYc8iOpA9qY1XHKNsP5V4eY5zdRdpCCJ55ZCy5SRlCNr6yRm5/bW5EtsXF6BJ7gaxxBDsIkXpANZ44l2xPKyDD2dSqQGle+/JWzt6l5G81ivLAmTih5U7PJIErTEff+q75DoCD3HIWKnBlkRreNwrgCaMawkpEWrlwDWRIdsGSXIlmPTfNjZCjNjL5uqN+pz+xFBkY+4TKy4rv4ApmXfVFUrY5QDaHYGgunx6BeJ408TY63GxYSqw0JfSTr0YYYshbWmnfa5g0UUjJ2x0JWqhos8CBe2Lr7TfqY/uNOgaPt9/uH/9g0hwO+Kr5dqeI3DdAmPuGKD+1HfMJnfOm6u3j7801ecN3F7fWMW1f3JA7eq+Ai0wTjspritX3GC+61Xo/N7wAAAP//cXedB20CAAA=")
	data["internal/docstring/def.md"] = decompress("H4sIAAAAAAAA/1ySzYrbQBCE73qKAl/sZSOTzRPsIYdAjoEcQg4tqeQZmB9luse28vRhxoaFXBrRP9VfteaA48IVSSKx5hJPmHxaFNJTuslMMFnZh58+hF6E4CqhEpafbfAJ5oi5lsJkH6OvuDk/O3jFuSrLGdOOhavUYCPen4I+XRRSCB9jNZkCIWlBodZgTVsSWEoumOjTBUW8coFfW0HMGDdrK6IsHarwU+c0xzjih/OKiU6uPpfWtvh1Zec0Jwkxq+G7101fIdqGdtya1wsTi4SwYxUfoD4wWdgbkNbZYRaljsNwOOA94etd4hY4DHhc9D4AwDHK1j+A45rwa/+NHS8veDs9s8o/n1t4a+HL6TQMnZcPtQfIXChGCIL83RFlw82xEJTZgYGxeclr/wXmComt5KtfuDTRyjRTu/Fcp8AFdcsJpVXURnyzxxZzTHjeTexjl+Uu/N9zwPl+Hod/AQAA//81Ap3HPgIAAA==")
	data["internal/docstring/defmacro.md"] = decompress("H4sIAAAAAAAA/1SQT27yMBDF9znFE0hfDF8F5QBVxaJn6AKxGOyJYsmxU8+k0NtXNoGCF6Px+735o1nCOO4Gsjkh0sA4lLheH9GlPPxf4eSjExAyk+OM6mw+fQiVgK4KTj/Xck3QnmGnnDlq1WQkyy8499728ILtJJy3pcJxR1PQDfZzFy/gy0jRsSu8dJrn+oiUS6IJFLQkPUPSlC3DJsetwJGWPcfMwlFJfYo4cZcyw2vt/U1hImW3aZrlEvuIjwsNY+CmwcMdbIquAYCFL3DgqFLHtYW09TKLajj8gw00Ccux/s255wgj/PV+A6sKAOM7mDfsYALHO7xRwMwSXh+04EXRUuCt7+7qk/cv362eHCnKdV+YzKJzvM+tr/kNAAD////E7O/+AQAA")
	data["internal/docstring/defn.md"] = decompress("H4sIAAAAAAAA/0xQsa4aMRDs/RUj0dzpAfdIGSlFiuQLIqV4esWeb82tYq8P2ycCXx/ZEMCNPTPe0exs0E3sFEqB8UHpmD/hYgpvPUbRKYOalBeyDLeqLRLV/Bbvmw56kBgvN5cSUWaGXVNiLc/xLc6z2BmSMayZ01AnJna0+rI3ZrPBd8WPvxQWz8bgHszJaADgQz7b3dmoU3sB3TcI3nu8v+JDj8Mr/vLEX9lnrvwbOicjul3T+xd06Osx5tcsGXwLg4VT7SQjsV1TrsuepczQiELiYcl7xKVIkCvVLrYgnaocOMQ7t8fPmEAIMT0cSUv7mQvZPzuXhHXyl7pzVLJWkPm0slrGkZVT8/mfaovM3Jqeol0Da7nJ0WHwdL3sMp+GvfkXAAD//5TJ3kPhAQAA")
	data["internal/docstring/do.md"] = decompress("H4sIAAAAAAAA/zzOQarDMAwE0L1OMTibJHz+HbroGboWjdIIFDnYcunxS1zISqA3DDNgXDLWXPZ5niBvtsYhFXuz0MOkU6WHml0K4efWAeqIVvwPRc6r/kJsglWdr7xmB1doVBSpzeKfaBhwc9w/vB8mRDhXEACMR1EPc6RNzHKa+jf96mVJE30DAAD//3SexluyAAAA")
	data["internal/docstring/drop.md"] = decompress("H4sIAAAAAAAA/1SQT0vEMBTE7/kUA3uwAemi+0evHjx4FzyILDF5pYE0qXmv2vrpJZt1dY95M/xmMis0LqcRNk1RwPShUd4M6QmdzyygQANFYaQOplgmipbUiw8BmWTKEQbBfC9nDdIbAc02TI7+o9bHmPUFsqhjTp/ekTsTWjxVhdOULf2RPYP7lIVyCYlHTy4iC7n6i2uYCBpGWRA8C75K0Xc6dSXXKrVa4SHicTbDGEgpoHHUYcZVc4NbbLDV+ve44HWHPe5w/1ZvZa0NGpuiNYIZi9ZKPfeeQZVXA0/LlH6X2xyaLXbY60OrfgIAAP//yGe8XX0BAAA=")
	data["internal/docstring/eq.md"] = decompress("H4sIAAAAAAAA/3SQsW7bQBBE+/uKkdRIiED16VKkCBAgZUphdRzqDj7uibdLy/57g6IMuHGzCwywM29HZaT9BCdsOIUd9pww1DY+xo8DnOaGPEBgdNQBr1JmGqQRuad6jlLgFZ6IITfz8D+Xgkafm+I8SDGeIQarVZedHdRYZ3U2g6xxnsSRDVr9G9sO/zyx3bMR968J3maeuxB2O/xS/H6T8VYYAgDsew5I2CaWUreHVeP0KSAdQvibX4ixmuPW2OcoTjvCUzYMs0bPVRFFIcUqLoTyKs4el/fl4Ebts14fmKfNCTFJk+hsHf7o6hLFeHxiruQXPuHZL9UuP0spjyLWYjnNUrrwEQAA//9n/ojWnQEAAA==")
	data["internal/docstring/eval.md"] = decompress("H4sIAAAAAAAA/1JW0EgtS8xRSMsvytVUSK0oSMxLKVZIzEtRAAmXJpakFiskgmW5AAEAAP//Q1NDtisAAAA=")
	data["internal/docstring/filter.md"] = decompress("H4sIAAAAAAAA/1yQMWvDMBCFd/2KB1lsKIUkHbuU0qF7txCM6pzrg8vJlU5JnF9fpFDTdLxPeu87boVmYDGKGLL2SPTdQvyVZcaNJ/hCM2lP7jWSNypI/HVeOM5jSIQ+qJEaOMFGQqSUxRAG+GmSmfWr4imGEx/oUIXGQWGhPpDQkdRSSdx9/NU84n3411ym0s69r1VFHbONMxoNhm7wkqh7QB2UpWtLRGvu5CUTziyCTwJrL7nYWP84ytKL3rnVCi+Kt4s/TkLOYTleMyh2lz2aZ1ywbVvs1thgi6d969zHyOnmiWQ53vrvD9g1a2zazv0EAAD//zyu2KSPAQAA")
	data["internal/docstring/first.md"] = decompress("H4sIAAAAAAAA/3yQP0+EQBTE+/0Uk1whJAYSvRipjIWFvf0dLLPy4rIL+8fDb28Ao7nGdibze2/mgMJIiAmRc4nAlIOLSAOxy7Qc6RK82cTIOdNpqrdBIkx2Ool3uIi1P+F/shO1GGH/S7mFD3BiIdd0SATHKX1VeE24+Gx7dERHRyNaWovkoQfqDxgf0OKztfKHRY7i3nEuIuenrdgZHY0PhG6tXb16+7Be79eBMdWVUocDnh1elnacLJUCip4GC26KpsHDEfd3OD6W5Wrs/ZZSKWxDcA9d7XBqmlOlvgMAAP//i7WKe2ABAAA=")
	data["internal/docstring/fn.md"] = decompress("H4sIAAAAAAAA/2yQwY7iMBBE7/6KEhw2YQVasiza24rD/sFIc4gi1HE6YMluR26bgb8fkXCbuXY/PVXVGtUoEAr8Dy2li3YYYwo/a9jElFlBApIojxCLYixis4ti3p33L+R7AvlKGYEe6BkTqfIASrHIACcgjC5p3lpPqggkwmlnzHqNk+D/ncLk2RigGnjEEEvv2QBA5TmjDcWj6eYD5vyrBVHcyBfWFdp7h2qzwR2h+Lqu61mW4/bGNse0uAJNLzd+VHs0+I0D/uA4429Xp+AlCT6ebRPnkp69GIsF57bBAUf8xf4X9k133hlz+rKEYkpxKJZBsD5qSbyMY+PkWGdh72RwctHlcaUbo2cWJB45sVgeMKYYZlZLmnd0coHaOPHOfAYAAP//zN9MUMYBAAA=")
	data["internal/docstring/future.md"] = decompress("H4sIAAAAAAAA/1yRUZLbIBBE/zlFl/0jqRLvGfKRC6T2AIuhtaKCBxcM2tXtUyDFWYfPGabfTPcZw1y1ZmJO+TZNI7jaWK2ywKJQkebeKrBlE7fkJKmWuJlf1Jql/ToEgkCXXacNtbo4DUkueF2Ie05r8PSH2keIEVc+cL7NN+LdZquELpnWNyF+0tWm8w1WPKxscDbGAk078MDsktM1Jvd7QhUN8bFQwWJX4koKXLrdI5Vx+we/GHM+44fg56dtXWOAwXNut/11yADAoOn7SqcpY3insO3aG8DAW1CcMv1pfC6lbOWd/1c3xpg+TmN7nTdXHY15XUIB9zWaab46fnG5nU6Pl7nqC3SxCpdkZdbSj80sNWrpAchTZOjGIIgmWOw37Mm8dfJblz5y6V+/WLiz8yPy5hsv5k8AAAD//3S7uL4+AgAA")
	data["internal/docstring/generate.md"] = decompress("H4sIAAAAAAAA/1SSTXLbMAyF9zrFm3hjpYlyhi5ygU6mm04ngUnI4pQCVP44UU/fISk5zlKA8N6HBx5wPLNwoMQYNczfeuzfEYTIfzOLYVBcxUxBRXP0a/d8IZ/rP2lixIWNGx3bKhHhpI4uVGXTFJgsdAR/sMnJqQz4wSkH+WKRJkp4d97DpQakFw4gWcts8Smm1ZISKDB4dimxHfCzNW5qyNHJGQSvhrxfEY0uBTCLKQS7ZAHG27FMNfn+bcDL1tgUKzVbUME1+hg0Jyf8gJlJiksDQjFC0orQFrk/eTV/7pElOV/9jIbAcVGxZZA9zywJLiJwVH9hi9NaXSTmmcOOuYc0dN3hgO+C5w+aF89dBxwtjzDqNcTPY3YA0Na6C2zv+tuCBpIzf62t7L2+3/V9XzWTPl7YJA2bct91L5OL4Obb1js5sRXP07/185Ch3ratUrrXB7YHVIpCM+OpiT+1yG88y4/NogR4zY/KkYKjk+eS2ZjLZbes7AOo4chG0N5MU3z9VWO47n5d+Pfr0P0PAAD//3asAtAFAwAA")
	data["internal/docstring/gensym.md"] = decompress("H4sIAAAAAAAA/0SQsQ7aQAyG93uKX8lQIgF9AIaqQ4fu3SoEhji5k5w7OPtS8vbVNYGOlv19tv8Wu5GjLhPU8rcO98xkrCCUGJ6Foct0S7JHUR6KIERMdM9J3c8BVKEQRwTFI6c59NzvYZ7s3fgTRHDjSvewhGchCcMC87wtkAUjR85k3G/Ljvjlg2Io8W4hxbdaN6rnLEuV39jTHFLGkDI8qT8YBflotB6rSzR6HXwYvYTRW+XWB47OtS2+R/x40fQQdg7YCRt+6yeTZqbcdGcH1F5QwxcS/rpNNZ5FUnOGdl2lT6caBD9LmEk4Giw54HpdrTPl9j9Sq865f4/yesAaVmYrOeKyTV6O7m8AAAD//44FE0yjAQAA")
	data["internal/docstring/get.md"] = decompress("H4sIAAAAAAAA/0SQvW7jMBCEez7FQC7OBg6u7ho3QQKkSBukN1biyCJik8ruyj8I8u6BGBkpuTM7/HZWWB/o6Iue8M4bInuZjv6wgdI18UyD4CzHiWhvsyW80ifNBh+4CJfkQ8oQGD8m5o7wQRzJIGalS+KM1VR3bGSX+sQ4p23x0tdp/bzQkIuD12R+j607S/Df+ZXB5AO1KgvwQpJs5p40M6IoJIOqRetckjFuQ1it8JjxfJXTeGQIwDqyh5aW6vjcZTkRzVNpG+zkQPz7/7WZTXNNi2knMSrN0Fy05EOzCeFtSAb+ZC4Ihv2i79Gyk8lYie+dpDPrxfmPoyvZZa4w/4aPWkaq37bhOwAA//+YoO0YpQEAAA==")
	data["internal/docstring/go.md"] = decompress("H4sIAAAAAAAA/zyNTY7rIBgE95yin72BLHKHLN4N5gJfTNugwXwWPxl8+5EnUralqu4ZdlOsWvbbzUHqmZdQNGuv6QRfkro0VgieSZdv8xWIo+grevq/quInpoQnP7JHzBBUHlKkES0Uioeu4ODSW9R8xyOfKKw9tZg3XCUvowW+jz6rPtZFiqe/GzPPeGT8H7IficYA1nPFgD2K7rHSuYttCjswBaak0xvVVmCHw4RT+7/Jmd8AAAD//0iWKE72AAAA")
	data["internal/docstring/if.md"] = decompress("H4sIAAAAAAAA/3SQu27rMAyGdz3Fj2Q4FmAco01vU4sOHfoUgRxRMQGaCiQ5cd6+kI02XToS5Mf/skXDAadEHmUgBUmmN4sTpRDTmJF5PAmhT04PA+vRfIZ6CDo7mVwhv7B8cIXAGSVNZbii0ViwD04y7Vssg7LsbbuwXVXqUAUqc3vl1CNRmZKSbxHLQOnCmVaoOluhFhzg9NriwiLo6Y8X/43ZbvGu+JhdTWEM0HgKmPGvucM9dnjAI57wjBdrly0HNK9ohBSzxc4aANjM1WbPx82vMY9OZGPNdyH5OvZR0E0qlHNXT6ZMHqy5kPOIAR2Hrl1rroTEIx+c3Lr9icN6plRqgK8AAAD//9J6mECfAQAA")
	data["internal/docstring/is-assoc.md"] = decompress("H4sIAAAAAAAA/4ySsY/VMAzG9/wV3/UtnDjeCdhuQQwMSIy3IfTkpm5jkSZV7PY4ofvfUZrqCRADS6Ik9ufPPyfRzPoA0Tekmj329QNu2u5OeHXcjLnMr29hrKZ4CmyBCywwlpI3GXjYIxRUuIkImWys7vMIivF45Y3iSsawDEot8K7KJFgQxbgmb5ITniRGFLa1JFysrHw54zEwRilqSDkdfo84XaNB0m7oqtHSJU24jBSVL2fnTid8TPj0g+YlsnO49vfzoaJA10uM3cv1uPA0PXcv+PoW7/D+261zj9UnN4E/bbYi6NnTqrx7sSClkYEoCBt7y+Xs3Bf5zpizGpbCg3gy1ru/GPhKKGpGz0g8kfGA/rkmLJyG2lctcX9zDx+okDcuFZIoZqaksEC2s8/jvyY1rztJq/K/T+y8Yzl+ALrAMebuvwC0OblfAQAA//+FUUXoVAIAAA==")
	data["internal/docstring/is-atom.md"] = decompress("H4sIAAAAAAAA/5RRu25UMRDt/RUn2YJdBTbiWaRBKego6RCKZu258Qg/rjzjXfL3yL4kFRQ0tjwzPnMehTLrHUTfkNWMcXzG1bzcDvvtvdSWbw4wVlNcIlvkBouMtdWzBA5zQkGNJ4J4dz9LEIWvRSVw4wCrOD1PQBaIwVMp1UZ56W3i8plSJ+MAKgGX2lNAHZ2LKL90B5aYclqOzu12uC/48ovymtg5PPMuknAXOaWK64HA1wfnvkVR8DaLi6SExtZbwYO1zg84sac+NpGPGNt4qNhIH537Kj8ZuaphbRzEk7G+hg3QpRdvUssQBUo61RZ+nGJOT+PDyiVIeZzm3V7dwkdq5I3bEZNXZioKi2SglFCXv9mcuxr+uPbCC9hvseHV/gZv8Q7vD/j+AR/x6cd/qPY1r7WXAHtaWaHdR5AiyYh+BHJmb7VtWQ8O/4j36H4HAAD///Tfnd5aAgAA")
	data["internal/docstring/is-indexed.md"] = decompress("H4sIAAAAAAAA/3ySz24UPRDE736KSvbwZfUlG/GfcEEcOCBxzA2hVe9MzY6Fpz24ezbh7ZFnsgsKgZMtq7tc9etWGWjvEO0qast7tng43+PseAsrXJxeu1yG/9dwmhvuenrPAu+JseRDbNnOFQYpPErB+H2iNrTwqYOk9FDCg6RJnPAMwUFS/FUK70XRiGJHSNPQjC12PxbNy/qjwvto6CZtPGbFXUwJhT4VxdbLxO0Gtz3RxWIOzXr12M+xxabkiDrHOMktSlH32HaSjNtNCKsVPig+3sswJoaA37j8d/EMz/ECL9f48gqv8QZvv65DuK0WuTQ85TCEz/EbMWRzjIVtbMRpl4+iVQ6SLFcYyr34wmIsHKlt9VitX59do+mlSOMsNXs0DBS1CtNn8Ll7albDNAPyKr+M4Y/Jbea0p5XAec+U8jlubv4ZciGHHRuZjDAvUfd/WY5N+BkAAP//tpn/zY0CAAA=")
	data["internal/docstring/is-len.md"] = decompress("H4sIAAAAAAAA/3ySzW4TQRCE7/MUlfhALBJH/BMuiAMHJI65IWS1d2u9I2Z7lukeB94ezTqAFEwuK61UXVP1datMtHeIdpWoSNT3OGvfsMLF8jfkMj1dw2luuBvpIwt8JOaSD7FnvygMUoguV3X2MH6v1I4WPg2QlO4lPEiq4oRnCA6S4l8pfBRFJ4rCORdHdGtx9j5etucUPkbDULXzmBV3MSUUei2KrZfK7Qa3IzHEYg7NevUwzO8Rq8kRdenwx+7oFHWP7SDJuN2EsFrhg+LjD5nmxBBwD+TJxTM8xwu8XOPLK7zGG7z9ug7htsXjUXwqXQif4zdiyuaYC/vYidMuH9RqACRZxo5Q7qUV2P1sAzO1b/la7Ouza3SjFOmcpfWOhomi1ij6QjwPp5Y01QWON/sj/39WtlmaLjeA85Ep5XPc3Dxa8EgMO3ZSjTAvUff/uYhN+BUAAP//S6TElXICAAA=")
	data["internal/docstring/is-list.md"] = decompress("H4sIAAAAAAAA/3xRTQsTMRC951e8tgdbrFv8tl7EgwfBY29Flml2tglmkyUz2+q/lyRVUMRLQsKbN+8j0sTyHl6eBS+KcnzAql5mg217jylPT3dQFhXcHavjDHWMOaebH3ioCAFlrgxiPo+gEB7ffKOwkDI0gSpgX6Yj1HnBuESrPkXcfQjIrEuO6DUv3Hc4Ocbosyhiik3iAyZLUPhYZfymaNM+XtGPFIT7zpjNBh8jPn2naQ5sDH6ZerJ9jhd4iVc7nF/jDd7i3dedMaeiiRv6T0mNERe2tAjXxcI2xeYeXkC4sdWUO2O++G+MKYlizjx4S8qy/8uwpQgKknBhRL6S8oDLjzIwcxyKi7LjsDrAOspklXNJxAsmpihQR1pjTuO/2piWGpsW+tpKV923crF2HEJa43w8/t92a8L8DAAA//+PHUBwKQIAAA==")
	data["internal/docstring/is-mapped.md"] = decompress("H4sIAAAAAAAA/4yRMe/TMBDFd3+K908XKkoqYOuCGBiQGLshFF2SS2zh2Jbv0lKhfnfkJC0CMTAlPt29e/d7gSaWE5y8mSgl7rF+PuBl+zE7vHrUhpin13soiwqultVyhlpGyvHieu6XDgFl3nTM5wHk/VbnC/mZlKERtHVAb4kPRSVArRMMc+jUxYCr8x6Zdc4BjeaZmxpnyxhcFkWI4WF5a5TZK1xYDD1F1nkXRjQDeeGmNma3w8eATz9oSp6Nwe8Df54KD1St8766P5+Jx/FW3fH1Ld7h/be9MefilFeFP42uW9ByR7PwYkatyysaOAHhwp3GXBvzxX1nTFEUKXPvOlKWw18UOgogLxEtI/BIyj3aWxlIHPpyWFlxfDmis5SpU84FkxNMTEGglnSJIA7/imqaF5Za5FcK9ULkET8qy97H6r9uX0P6FQAA//+twuBbUwIAAA==")
	data["internal/docstring/is-nil.md"] = decompress("H4sIAAAAAAAA/3yRT6/TMBDE7/4U814PvIpHKv5cygVx4IDEsTeEom2yqVc468i7aeHbIzsVEghxiRR7Z/Y3Y6WZ7T3EXqkkqKQPeKjfsMNT+5tymV/u4WxuuEX2yAUeGUvJVxl5bBMGKlzl4fMESul+yFdKKznDc718rkKFRzFMqw4uWXGTlFDY16LovazcdzhFxiTFHJq1kd2nbE0O0Qbw22ETi17QT5SM+y6E3Q4fFZ9+0LwkDgH3NC+eXuMN3uLdvvLsQzhVFt7G/kTZrHDmgVbjbWNDqtEgBkIS8y6EL/KdMWdzLIVHGcjZnv9KOZCCkmWcGcoXch5x/lkFC+tY2euCw8MBQ6RCg3OpNYhhZlKDR/JWbJ7+1f68tq682UvqWuL2kHiMnFJ+xNfj8dt/A2/dh18BAAD//8PWqwgSAgAA")
	data["internal/docstring/is-promise.md"] = decompress("H4sIAAAAAAAA/4xSu47bMBDs+RVju7EQR4adLk2QIkWAlO6CQFhLK5MIRQrclX339wdKNoEzrrhKwj5m58FAA8t3OPk6pjg4Ydy/P7B6/JkNtqXaxzR8qaAsKrhZVssJaue9q+u4mycElAqUmN89yPt7h6/kJ1KGRtBjZpcxAtQ6QT+FVl0MuDnvkVinFNBomripcbKM3iVRhBgK6fukTF7hwsynoCwALlzQ9OSFm9qYzQY/A3690DB6NgbYdtxjPBShVVWKx6disWI85ObfA4749q8y5pTJ84L5nvtyF2duaRKe6al1afEKTkC4cqsx1cb8cf8ZQxTFmLhzLSnL7smYlgLIS8SZEfhCyh3Or3lh5NBlqfnEfrVHaylRq5yyc04wMAWBWtI5kdh/lN0wzfZqhn9kWM8ulTeBtWXv4/pT8pfozFsAAAD///2Gr6RsAgAA")
	data["internal/docstring/is-seq.md"] = decompress("H4sIAAAAAAAA/4yRO2srMRSEe/2KsV1cm+sH95GH04QUKQIp3YWwHO+etUS00lrnrJ38+6D1AxJcpBEIzQzzjQI1LHdwMhPeQnh7j0E+zQjj/lbH1PyeQFlUsLeslhPUMtoUd67iqlcIKHG2dxxKFvNUg7w/PvGOfEfK0AjCjryrztJpzgpQ6wR1F0p1MWDvvEdi7VJAoanjYo6VZdQuiSLEMDvZT1LpvMKFvtg55pDgwgZFTV64mBszGuEh4PGdmtazMThS/hr/wV/8w/8JXq5wjRvcvk6MWeVafBBfamXMs3tjNFEUbeLKlaQs0284JQWQl4g1I/CGlCusP7Kh5VDlfrn2YrBAaSlRqZwyrxM0TEGglrSfM9aXlm+6fhTN8V/HlXlP2H8ohpa9j0Mslz8B+wwAAP//P8BozhgCAAA=")
	data["internal/docstring/is-str.md"] = decompress("H4sIAAAAAAAA/3yRT2vjMBTE7/oUk+SwMbtx2D+X7GXZQw+FHnMrxbzYz5GoLBm956T99kVyUlIovRiMZkbzGwUaWP7CyUY0QTT9wyJ/zQrr8tfHNHyvoCwqOFtWywlqGWOKJ9dxVxQCSpztLhzF3Pcg7y8HfCI/kTI0XgU/ckCAWifop9CqiwFn5z0S65QCGk0TNzX2ltG7JIoQw2Z2X4UyeYULpct7yOzPqqYnL9zUxqxW+B9w90LD6NkYXMC+rX/iF37jT4WlZe/jsjJmnyvxLP3YaI7DgVuahOdbS7NMCScgeCdaG/PgnhlDFMWYuHMtKRfiW9iWAshLxIER+EjKHQ6v2TBy6HL/fMF2sUVrKVGrnPIaTjAwBYFa0rJx7D97jGEqk2mOv2xeF/LFLXqFx93u6Uvq+R3MWwAAAP//H5F1yCYCAAA=")
	data["internal/docstring/is-vector.md"] = decompress("H4sIAAAAAAAA/5SRva7bMAyFdz3FSTI0QVMH/W+6FB06FOiYLSgMRqYjobJkiHTSvn0h27l/uMudJEiHh+cjI3UsX+HlzYWtpozp+IbFfDErrG9vbcrd6w2URQVXx+o4Qx2jz+niG25GhYAyzz5ifragEOYPvlAYSBmaQLNkWxwi1HlBO0SrPkVcfQjIrEOOqDUPXFc4OEbrsyhiire4s1CGoPBxDHNnMtX7eEbdUhCuK2NWK3yP+PGXuj6wMbiHe7V+i3d4jw8bHD/iEz7jy++NMYeSiyf941iTJ05saRCeWo/xCiu8gBC8aGXML/+H0SVR9Jkbb0lZtk+ILUVQkIQTI/KZlBuc/pWCnmNTIEqD3WIH6yiTVc5lJF7QMUWBOtJx0ql9biXdMM5Ni/28mmrEv+0ZS8chpCWO+/2LuYVtis0D8Mm0Mv8DAAD//zZhBpdfAgAA")
	data["internal/docstring/last.md"] = decompress("H4sIAAAAAAAA/3SPsUrEQBRF+/mKC1uYgGRBg+xWYmFhb787mdwxD19mkszEjX8viaJsYfsu5/DODoXalJE4lpiY5ykk5I7YrlT2DBnRb7fEcWZwNK+dJPg5uCwx4CKqP+z/6EAnXtj+Sm4RJwRRyLUcksB+yJ8VXjIucdYWDdEw0IsTq8gRrqN7h48TLD6syp8Wc5LwhnOROD5uWWc09HEinFVdt/364L4yZrfDU8DzYvtBaQxQtPRYcFMcj3iocX+H+lCW67A1LaUx2Nr5zVyln+rDqTJfAQAA//8Cb6VGUQEAAA==")
	data["internal/docstring/lazy-seq.md"] = decompress("H4sIAAAAAAAA/0yPsW7EIBBEe75ipGuWiy5ylC5divxBOivFghcZBYMPcOS7r4/AOiXTPt4wewIFvt8uRa5wKS/ns8aa07RZKWAUuW4SraDOXOEL5IfDxlUmBL77cFPqdMJ7xMfOyxrkTSmAJnERzpveOn4pAKAgFaPzBnTAkWEO9D9/a8im2CZQcwzoCQyjWx5WJwNetNb9W5fyRdjOGHdQ5W/B6zD0V63w4dGafawhYm/a59yOOsajk4I6C5zPpXbfeZMiW+sRt8VILs/qNwAA//+YjIqANQEAAA==")
	data["internal/docstring/len.md"] = decompress("H4sIAAAAAAAA/1SPwWrzMBCE73qKAR/+GILhb/MCPfTQe6HHIKvjSCDvOtKKNH36YjtQelx2vm93OhwyBZXXHoXWilRYJDLlYhE6bVPltVEC3UfK+ZHbFtLmkWWNMXOmWEUSeARtYn7Mv+iAt78upIrsv1O+I+i8NOPnEb7eJcSioq0eoQVqkeWWKpEk+GVT6oSRSS77lRWzmCqmJsGSCm7bk36FvIClaBmc6zq8CF6//LxkOoe9+L/DfzzhGae+d+591XBPPCx71fPpPLifAAAA///Twvx/KwEAAA==")
	data["internal/docstring/let.md"] = decompress("H4sIAAAAAAAA/1yQP2/jMAzFd32KB2S4GDnYuP83FAU6dOjSqUCHIIMiM7UAhUxNOrb76QvLTtBWiyCK78fHt8I6kWF7w/5IOEh7vN3s8r0psI9cK5IEn3D2qSN1zzElhJa8ETyY+uVbg5zoe1ZEfoE1BD1RiIdI9aKFCazxNvdiP2KaWeLB0E9Ua4hBU+8E/0yYDCn6aE3kjxDPNVqyruUsaEm7ZJBDfiWvdgFG4dK51Qp3jPvBH0+JnMOy/IBv6x/4iV/4XTgsZ8T2D/7iH/7vdrm4DsLBGwaMReHcUxMVNKPmBa6xpKiG4FOiGtVQZZceZwom7bU+Vl/dJ/82Yh5CnC3Pm4gSlF474kBa4lFyPN6yyHrJOSpqAYuBhmm6dKaxpksUVSKrcoqlew8AAP//tr4OOfMBAAA=")
	data["internal/docstring/list.md"] = decompress("H4sIAAAAAAAA/0zMO2rFMBCF4V6rONjN9SVkDymyhtQiOkaCkWRG49fug1+Q9p8zX4+XpGYYq+b3e8Cv0hsbPApXHCf3k0Tu/i9jjbURFGYWa/BKWCS4eJm9MZxkw6R1SYHhA1WhtFnLtcuT7ZeURpR6zw/mefl0ru/xVfC9+TwJnQNegSM2dJEitRuesqOzSOVVTnXDPri/AAAA///Sf+9X3wAAAA==")
	data["internal/docstring/map.md"] = decompress("H4sIAAAAAAAA/4ySP2/bMBDFd32KB3iIlBSCLSdB1qLo0KFbtiAIrtIpIsB/5ZGx1U9fkIple8skgCf+3o8Pt0FtyGNMtofw37sGmv4pPcOQl3yS2PYs1Y/AFFlAeT6vAxwmJwzWbNhGAQVGnBiBJekIN4K817Oy7+XYB/ehBh5KXlTOIroyWHknUotfI4wrOLJw9uIfJSvoW76twrVA4BgUf/AAZeEpkNasc5KkLAMaBpXDSYPCe1oufooY8v7Cr8Vv8j7rH5TWiByMshQZJBDnbP6Sna/c+DhRkshDW1WbDb5b/DyS8ZqrCkvb9WjxcnxFfXuLI7qmwcsOHfa4f22q6nlSssQFjinY4nXd+lvd4R6PeGreWjxPjNFp7Q7Zk5cseA6jC0bKi/Lgoot2Nbk7J+OmXpC4qfdbdFvstnh4bL5qtH9A94Rul5VOpTmrZ7i+T0EQp8CMqAxLdolTqWqR/cM9JeHPVeidHc7kwpjoBFgXpPofAAD//3FiDwi7AgAA")
	data["internal/docstring/not.md"] = decompress("H4sIAAAAAAAA/1yOQWrDQAxF93OKT7xpoCRn6KI3KHQ5TBK5I1BGQdI4ze2L7YLbbKX//n8DXpoGRrXrHqJffC4iD3CbyMIRlRDWo3Ijd+i4XG6mE1/osmDpk0VgFN0a8ljEKYOfgsc5eQT72vZ4hUYlu7MT7n/5sE75kNIw4K3h/btcb0IpYdXcVRLR3T6lj8r+j5znTqpCZbM40bl0p+U3FemE/NuQN5VD+gkAAP//IPscxwQBAAA=")
	data["internal/docstring/nth.md"] = decompress("H4sIAAAAAAAA/3SSMY/bMAyFd/2Kh8vQBLjm0Dadixs6BN3a7gfaomMCMuWKku/SX19IdjL1RvE9kY+ftMNe84ghpgmint/geaAS8rcDEuckvLCBsFAojO66mtxPziWpIY+8SXmkjJ4UHWOIRT0oN9lm7mUQ9lv/OECywfhPYe35iPPQfKsqhlhyNdVaVxvZ7XS78lhPCpY8cmrKlnmLIlajl6TsERNIwSnF1Ookxv6IH8wzRDGJ+jV6G19XTYy/nOLHjqrTud0Oz4rvbzTNgZ0D9p4HBHzYf8JnfMHpcKjFSjHghIfXFPXycHDu9ygGXu9tgQwvm/6Cjnsqdtv7hH1d5Gse18Khpu34Gls+Bqu/Y7gDDWJ5TXhu6T1+bYjwPM9BesoS1bmzWmZqDYqJXlqbJ83jE4aifTU9bgD8nbLBSj+CrI0xkHos3OeYrL0zBYsYaWFQupSJtXrqVPbwkrjP4Yoc66zp+D63/zB7lRBwFQ7r6kZT5Wf1gWn9cnPiRWK54z26fwEAAP//JdYY4cgCAAA=")
	data["internal/docstring/or.md"] = decompress("H4sIAAAAAAAA/3SPsW7rMAxFd33FRbIkeXkO2gIdumXI3LGjINh0JEAWDYqO078vZDuIO1STQNzLc7jFjgUtS3c47NGTlG+GQ/Ys+r8OUg9BQ7rixHIyl5uLg1PKUE+Ys61wh0itQhkSrl4rnDMyc4LL4ESgZ43hoDKo/0YZ0hFjiBFCOkiCeqfzvMKnepIxZELQOdQL10RNWbJsLGLFJNFdJ53KmO0W54TL3XV9JGMwXbj7hxe84m1vML3WxUzLf5NYFwFqNntjvlZG9t0ekehGsmbaqW/hUgP7u2+PGD0JufzxgK9hKcQH9i9kOSirTJxnyFbmJwAA///Ke/TyrQEAAA==")
	data["internal/docstring/partial.md"] = decompress("H4sIAAAAAAAA/0yQy2oDMQxF9/6KC9lkSBsaKIUsu2g/oBS6TDSx4hF1ZONHpvn7Ms5za52je60Z5pFSEfLYV92Bklt0cKycqHAGteciQdFTZougIFwVitGfzBeXmnRilUd8XvlxCJkhKmc0uXpgLRmUGDHxcx+qWpSA0sCYwlEs2yV+BlaUgcp9l2SIHsMv2yeQnm7ww9pRvEeWQ/Qn9DxVY52IFvCYeFd63ofE2JH3oq5hIYkTJX+LXhozm+Fd8fFHh+jZGGBueY/oa1693M+3wCveum4aX0brzpjvQTL4rJ4rpnYtbFbrDShPn69uwHbefKy7LUa+lGK7NP8BAAD//+QMK6ahAQAA")
	data["internal/docstring/partition.md"] = decompress("H4sIAAAAAAAA/0yPzU7DMBCE736KkXKxJWia0tJyQhx4AyQOCFVpsmlXctbGXvPz9ihBLT3Zmm/0aaeCjW1SVg6CLhRRZKX4iEwfDheU0U5JIenIvLL3/+iKgEUDjimUmBEG1LOwBnkaSTTfgKVL85/liMMP9ESQMh4oTf1zDz0NLNSDBfV0Tg0b0kXHwznlDAmKmMIn99S7hTFVhSfB83c7Rk/GAFbDrees10NXuMNbMz9rbHCPLXZ4QLN8d86YlxNn0J8BX9PYRFqSYG9tg5WDXWPjYLfYOdhm6dx+YX4DAAD///gpY+NJAQAA")
	data["internal/docstring/promise.md"] = decompress("H4sIAAAAAAAA/1SQMW7DMAxFd53iI1mSIMgdMnTo2hswEh0TpUlBluz29oXgxkVHiR/vP/KIUy4+ycxn5OKpRZ5BhmaFZ9eFExbSxuGDays2g/CbhxjqyBi8TPABhKFZrOJ2w/sAylmFE1apo7famVSebWKrV9RR5j2OVVRxeajHz8sVK0kVe3YuaCtHdTwYiVUWLl1JCIRIqn3ULV5SO7SOVCEWtaVto1f7Dfc9Hcngpt//6W6RbyEcj7gb3r5oysohAKfEA/Lfwc79L+Mwsqofttc5/AQAAP//etNsWFIBAAA=")
	data["internal/docstring/quote.md"] = decompress("H4sIAAAAAAAA/0zPQU7DMBCF4b1P8aQu2kgoVYELsGDBgh0H6DSZkJHssfFM2ub2KAkgtrbs9387HL6m7Iwh19Sgsk9VDT4yrHAng3C/3kEUPTkh5Z7DO5OKfsJHckQxN5D2sDldcjTcJEZodlwYfKU4kXPf4mMUQ6KuZoiBh4E7lyvHedujxCBDqVxY++V7UvC9VDaTrLiJj8sRlWxecxkZh+P+2LQh7HZ4UbzeKZXIIeBXdTjhEU94bpoQ1vW1bFOuo1GcK8XVgEo+cl1QCq/zCsygUn4KdUoXrjg9LJkyQBw3rgzCMGnnkrUF3nzB2axZ55Qn27KX5/8o5/1f2bkN3wEAAP//N4kNS4UBAAA=")
	data["internal/docstring/range.md"] = decompress("H4sIAAAAAAAA/1yQQWrrQBBE93OKAm8k0NdI/ugA5vMXOUEWIYi21LKHzPQoMy0j5/Qhim1IllVF16NrhyKRnBjBCQKtcDKUGBKTcgZhC82/h/b0cUXm94VlYOiZFHPizKIZembIEo6cMqYUA2xwYlE4GfyS3YVLaIQNtFoUvN7NCscrrJMhcWBRW+PgPSidlrDVUmLEWV0U8vWtc+SJFv/FjOibvrq1/rCdTH0FkhE2K8+/0rav8TQhir8iCj94cBlzihc38ljBbVq390dQ/ubUxux2OAj+rxRmz8YAQKHxz4UHjQmF0huju2/bNnAyoSvL0phn5z0S65JkW+x20r+0DdoO+wb7Dn+b1742nwEAAP//JVxzAJ0BAAA=")
	data["internal/docstring/read.md"] = decompress("H4sIAAAAAAAA/xTKuw2AMBAD0J4pLNHALixxil2kSaL7CMZH6V7xTlwuIyL9xlbA7cXTY6FNCn3kRA19ywZF0NL2rpbliuMPAAD//w7hDINBAAAA")
	data["internal/docstring/reduce.md"] = decompress("H4sIAAAAAAAA/1yRsW7zMAyEdz3FAVkS/EHwZ2jXokOH7tkDwr7YQhXJpSgnj19IdlKgk0HzjvcdtMFW2ZeOuJTYYZbwhszvHZa/GVLHwtjRfRpVjBlppraFIV2e+/3i8XGAjfQKBl4ZLcNSVfs4BEKZS7AqmiUUHnAal2zzKWLSNPuePa4lG0y+CLsliA6lnVrlXrNBYo/MLrXPgvAbydilEo3KHqKsRPDRm5ewBGfINAXPvtLZKPaEaBlKuRh134wLc+1ap0k5+1QyOgldCdLAfUbJNSs3zUq4Yu9xG31YICLv9sD861rbPGwH5zYbvEd83OU6BTqH52v9wwu2KnEgjjgedzvnTqPPuPkQoLSisZ1sXXF+/X8+uJ8AAAD//6X5wZLsAQAA")
	data["internal/docstring/repl-cls.md"] = decompress("H4sIAAAAAAAA/1JW0EjOKdZUSM5JTSwqVijJSFUoTi5KTc3jcgaJgAWSS4uKUvNKFJLz80pAdH4akjo9hZCMzGKFzGKFRIUg1wAf3fy8nEqFtNK85JLM/Dw9LkAAAAD//yVoAgJhAAAA")
	data["internal/docstring/repl-doc.md"] = decompress("H4sIAAAAAAAA/1SOQQrCMBBF93OKD91Y0N5BsDsXIl5gSCZ0oJ0EJ4Xm9tKKC7fv8x+vwynmgJTfS4+oXmZujpjDuohVrpqNbj9cJ/mf9t9BvUjQpBIP0xmawNYgm3r1Aa9JHepgPMfH/ZJtbkirhd0xEHUdroZx46XMQoRvE3uz0NMnAAD//xQLeSyiAAAA")
	data["internal/docstring/repl-help.md"] = decompress("H4sIAAAAAAAA/2SOS27jMBBE9zpFLWVh4APMbuDxbhaTIAcwJZUkAs0mwyb9uX0gOTFipLf1XnV17UJJuw73++stibuhLN6wJgg0czObrh3jgCnmsMFf4BiHGqjFFR91j+PVhST8jdOGz3F3arq2GqH2+eSwOJ2JoeZMLVAXaMkNfJJXY4p3exB77MNB6DLKQtiQSW269r368gBeqi9b/Hr8/69p3hYi5RhSQU+JF/iRWvzkaRv1Y8UvTFEkXjiiv22I1tAzI07gNWWa+air7AouXgSLOxM9qeDZSXWFI7w+ta9bYHd1jz8wr7PwWx8Gp7DkFKFK8UkI8UrbNx8BAAD//3pr7pSgAQAA")
	data["internal/docstring/repl-quit.md"] = decompress("H4sIAAAAAAAA/1JW0CgszSzRVACRxQolGakKQa4BPlyeubmpKZmJJak5lQqpFTApxxyItJ5CSEZmsUJmsUIimK+bn5dTqZBWmpdckpmfp8cFCAAA//9MhILkVgAAAA==")
	data["internal/docstring/repl-use.md"] = decompress("H4sIAAAAAAAA/0TNMQrCQBCF4X5O8WAbU5g7iKSzEPECwzJxF5I3IbMLensLhbT/X3wJpx4GxoBclC8LUFeLTbPJ9V9aMeS+78Z23BHPUgM1oHhM99vZuXwwd+ZWnaNISrgQ01vXbTER/KTZfZBvAAAA//9YAeiGdgAAAA==")
	data["internal/docstring/rest.md"] = decompress("H4sIAAAAAAAA/1SQvU7DMBSFdz/FkTqQLIkEFaITYmBgZ28d55hY3NqNf2geH8VFqdgsn/sd3e/u0ESmjMS5RWQu0Sfkiai/wdZ34lzoDdXn5BJs8Sa74HF1In8M9DaEPOkMLkbKyFuXdTFlUHimv5deaJx1HDeyw0fGNRQZMRADPa0zTgtygJlovmFDhMaPFnenUJLzXzg1ifNr9ThhoA2RMFpkzfq6QI8Q0a9efafUboc3j/dFny9CpYBmpMWCh+ZwwPMeT4/Yv7TtGtRTLK1SqP68Mf/0j82GHDv1GwAA//+AYLxzVQEAAA==")
	data["internal/docstring/seq.md"] = decompress("H4sIAAAAAAAA/2yPwWqEMBiE73mKIXtoAtV36KGH3gs9iEhWRwzExE3+de3bF5UtPfQyEPLNfPwXmMIbxpRnCyfCeZECSehTXJnl+NnfDoW3O2NP9eVDeLJ/UZmIJafVDxz+6cGP8AJf4ovAhUw3fNf4GI/egfcuxiS48jnJ4RVd9KHDY3deiUy558ihVupywVvE++bmJVApwDwmxipQ0JTzLD0xhKRtqwDASKpW9pIyzOwWmDGi2VqYIhkbdKWtRbHWKvU5+QKe06f7FKNr9FRpaO4RfiNVuu1q9RMAAP//+bJX0E4BAAA=")
	data["internal/docstring/some-first.md"] = decompress("H4sIAAAAAAAA/xzMy63CMBBG4X2q+KW7uUSYdOAK6CEy1oBHxB40Dx7dI7I6iyN9f/g36ZQy6P1QXEW7zfMBKaOXqoIXe4M1UU+VtQY7j9t05jthSXk54hIOJQ8dhnXwtqIYTGT8WsZnN0HPskVxsv1GbafpGwAA///BbWHeewAAAA==")
	data["internal/docstring/some-last.md"] = decompress("H4sIAAAAAAAA/xyMXQrCMBAG33uKD3zRYu0NcgLvUGJYzWKTlf3x5/bSPg3MwBxwNGk0pQT6vhR30WbjeMJmWi4q+LBXWBX1qbCWYOf+GK78JMxTSvMZt3AoeWg3LJ3XBdlgIn1j7r/9CnrnNbKT7TVKvQz/AAAA//8zB9jPfgAAAA==")
	data["internal/docstring/str.md"] = decompress("H4sIAAAAAAAA/5yPTU7DMBCF9z7Fa7uhEarEzwUqxAWAHULtJHkhlhy7Gk9SentUp1RI7FjO+H1v/K1wk03RJR2qao0mxYlquSwyLEGQTX38dE9KMWYIIo+XJTpNA6znZfadZ4tJwsiM1JWXg6bJt2znyo1zqxW2Ec9fMhwCncP8g2XPENIS73e4xwMeP9bOvfU+g3MQRx8ClDZq/HURu8Jdqd184IXSUvFaMtm5bTBqFPMTw+kWKRKNRIyZpWufTRd7dGNszKd49j5oaseGV/+LGjWXQC+GapBThZqQOvDM1IRSWtSn0roN80zd/Ggu/u/5B9xt3HcAAAD//2Omoji+AQAA")
	data["internal/docstring/sym.md"] = decompress("H4sIAAAAAAAA/3yOQU7DQAxF9z7FV7OgWZA7AOIG7CuTcZmRHE/qcQq9PZoAEquu7K/3bb0Bx3Zb0MJHzNWu4tHAPRf7QLGoYOtT3CSh3Zb3qvTy14wsWL1eS+rw7tGE5xoZl421nIsksCVonVl/Cw3sAp5nWUPSRDQMeDK8fvGyqhABxyRnZFGtj116Nz/s+TCOncvlH37Y15HoLZcG+XmDz6IKl9jccArf5DTRdwAAAP//mE7l2gUBAAA=")
	data["internal/docstring/take.md"] = decompress("H4sIAAAAAAAA/1SQz07zMBDE736KkXL4YulTKv40cOXAgTsSB4QqN9k0Fo5d1mua8PTISZuKo2dn57eeAqWYT0ITkhdE+tLI7wjpCZ3lKCBHA3mJCB1MtiTyDak36xyYJLGHgTM/0zrLTrLSE2MzB28QGB2diK9pHYdhpuzpYL23/pDXsnDk8G1bate8Ci/LJIbEDV05NiL2gYUY0hs/ezgPo1C7/On/WYzJSWasu6d8/54gnHxjhNpKqaLAk8fzaIajI6WAsqUOI/6VN7jFHe61vogT3reo8YDHj1mbW6xRNiGnYcSktVKvvY2gJW8hnhvLR/3tbHdhYIta7yr1GwAA///z6h7tmwEAAA==")
	data["internal/docstring/thread-first.md"] = decompress("H4sIAAAAAAAA/3SRPW7cMBCFe53iAVtEkhVpY8RtABcucoKU9qw0EglQMwI52p+cPiA3WSCFS4IzH9/3eED99Qf4ukXMGtfUtg3MRaYp4Uxh53zSfXEYKYQESjDHPmL2MRkoLvvKYtVbHibjhCHDBpBMD5C3B8UcI+3bFjxP9xc/Q/Z4lVsZgTky+ARRA4XMvIEw7zKaVynJcPEh4MQYVc4cjSd4MYUK48SzRv6bxsvSV9XhgFfB25XWLXBVobRwRP2Eb8cGddviuUE9DHhpmqr6ldl83YqT4iNf5Jn6Ccey8NzgpfkozpFtj4L37+89fgo0Thzz0onNOOLs007B/2ZcHNmXhEW9LFDpStTVL84wOtXEectL4mggDN2QiyJsgUZ2GjJ21lgKvZvxdP+x/j+f7mHU/XPqitWfAAAA//+mkkxa+gEAAA==")
	data["internal/docstring/thread-last.md"] = decompress("H4sIAAAAAAAA/3SRsW7bMBCGdz7FD3io6bqSGiBrgAwZ+gQdk7N0FglQdwJ5cpw+fUFJ8FCgmwTeffy/nwccf7y8gO9zxlXzVE4nDwuZaSi4UVq4/ukyBvSUUgEVWOCYkagYKI/LxGLurc6ScUFbWS1Ihgcn2gNigVGWeU6Rh+3C/xAbvMrXOgELZIgFogZKFfkFwnWR3qLKmgufMSVcGL3KjbPxgCimUGFc+KqZ9zBRxsa5wwGvgrc7TXNi57CV0OH4HT87j+PphCePY9viuem8d+53xfN9Xq0UH/vRNrmtofPef6zemW3JgveueXp+b/BLoHngXDcvbMYZt1gWSvEP4zOQfSsYNcoIlfMaeYpjMPRBtXDdilI4Gwjtua19EeZEPQdNFXvVvPa6GfKwvVvzrxfOu9n6sQucvXd/AwAA//8gW007BAIAAA==")
	data["internal/docstring/to-assoc.md"] = decompress("H4sIAAAAAAAA/2SQzU7rQAyF9/MUR+niNrolRWKXHQseAbFALMyM01jMT5nxNEmfHiWIComVZfl8ny3vsNd0R6Uki8Kf/1vYFC+ctaxt5Wi5QBMoYgsJqVwYRXO1WjObF/F+ZSwpR1IGobAiDb94iX8MHZ6jlw+GjquNoqPscPwWHTHUaFVSPGAaxY6QAk9X8QtsCueq7A4bmblUv23TUcoNw7Re9c4IpJyFvFzZQUJgJ6Tsl86Y3Q6PEU8zhbNnY4C94wEzXvtIgdGQ5wY9nRj33cNb+xNY8G/fTyynUdH4tTTtNry9ccbSmq8AAAD//200PqRYAQAA")
	data["internal/docstring/to-list.md"] = decompress("H4sIAAAAAAAA/0yQTU7DQAyF93OKJ2XBjFpaAb0AC46AWFRdmBlHsZifkHFo0tOjBIhY+vn7bMsNrJb7KFVR+XPn4Ev+4kHrUo6cPVdoAWFBzJvEuBCelDMpg1BZUdp/tOSNP+A1R/lgaMeoSjnQEHD8GXBEO2avUvIe1058B6mIdJM4w5fUj8phv5oD1zGuW7STumm4Lte8MxIpD0JRbhwgKXEQUo7zwZimwXPGy0Spj2wMYAO3mHB+wCOecLq4v2w2AGAT9bBtxnm6wO4w4eTc2rmzv4pbg+1rE2ZnvgMAAP//e2v3tEYBAAA=")
	data["internal/docstring/to-vector.md"] = decompress("H4sIAAAAAAAA/0yQwU7DMAyG73mKX9qBdIxNDF6AA4+AOCAOJnFVi8Qpjbt1e3rUToIe7f/7bMsbeCsPJw5WBlT+uW8Qip54sDqXI2vgCisg3CD3LinNTCBjJWMQKhtKu+JFV8Yeb5rkm2EdoxpppCHicBtxQDtqMCm6w7mT0EEqEl0lXRBK7kfjuFvMgeuYlj3WSf3TcJ7v+WJkMh6Eklw5QnLmKGScLnvnNhu8KF4nyn1i5wAfucXkAMBn6uFbxcf0Cb/dYsKxaZbozj/iiCc8N0tj9aj/BFPjfgMAAP//x3q+iEQBAAA=")
	data["internal/docstring/try.md"] = decompress("H4sIAAAAAAAA/2xSy24bMQy86yumzqG2EXv/oEEPPfReoIciSFQttyKgFRciN67/vtC+Eie9isN5aHiHvZUrOin98YjgLcRTSH5UOh7RcfYpXZeHhwOiz20iRfGs1IJKkaLuJ6cEevFp9EawSBiKvHBL7cSr9+C+p5a9UbpCoxQ7BS5hZOP8B9zB55kLvHLfwxsukUPEIJwNQbIVSRUxeK3qnE0mtY6LGprJfIPZLC5RlNAMhVoO3qjZHCpM8GRlpKczvnfwaJag2zIrWuo4Vx+r8KWm/E036mzwnVGBh44hkGo3bl/BkiHdf/4DUra1d65ZN5vt2bnmbSHNG1vwiuf9NMWv7PsqsQR9XMo8PLvmtsCPBMv8dcPd3eFrxre/vh8SOYfpPBwA7KdisIuUknzaHebHoXC2lLG7SP5sKORDRE/reHUIzt3D4wd4P40r/BavUCu3eK/1k63Uizl92UHXlTXDK7aVXBkPzv2IrHNx03C+FklJLpVmOZ8gWSXRuaZ9J4M5rAMqp/sXAAD//74dn7EtAwAA")
	data["internal/docstring/vector.md"] = decompress("H4sIAAAAAAAA/1SOQYrrMBBE9zpFkWziwM8d/mJuMDDrHqmEBHLLqDuOffsBx7OYbb3i8a64rYzeB3If8/0+IQ6K0yBQvvCG4au2dpI/AK/SjWDjTHWDDMILwVXaU5zp0BqW0deamB74LNWQnxq9dkU1aEeqOXNQHV5ED8Gpb9U5pMF2ddnALXI5Xo7qiKL4JvzoSqgKQa7D/F9sYoYsVmrXRwjXK/4rPjaZl8YQgFtixoZLYWv9Mv0uOy5eOPhezogN+xR+AgAA//+vfXdiKQEAAA==")
	data["internal/docstring/when.md"] = decompress("H4sIAAAAAAAA/1SQzW7qMBCF936KI1jcBHFB99K/VaUuuugTdImceEJGHcbIdkry9pUdSunKGn3njDXfEtW5J8UpkEPnw3G1qtF6dZzYqxWZQJ9WBpsowqIR336Ytw6ppytwpc2tTQSOSGFI/YRKfcK+sxJpv0YZlGVfr0s3fxVhw82WDd5Z5DqDbNuXHFiRhqBrBMov62HewWqvefYKG8EpIlAcJG2MWS7xongd7fEkZAxQOeow4k/1D/+xwx3u8YBHPNV1ocVE9YxKSDHW2NUGAKpTYE2iWIz5voYPiwuYc7X5NhKnY+MF27zor/q0zfkhkgNrTGQdfDfTbdGgpfXL3o9Uq67g4hzn7Ka5le5VJnB3OTd3s+TZvvkKAAD//6VYmc3YAQAA")
}
