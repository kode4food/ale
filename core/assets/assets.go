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

var data = make(map[string][]byte, 14)

func init() {
	data["core/00_builtins.scm"] = decompress("H4sIAAAAAAAA/3SUwZLTMAyG732K7I0u42G4EpZ3cRwlK3BkV3JCw9MzLW2HStrj//2ObP3SpO/7vosZulQYvnXDirkhHQ6fRpjCTXXh+Kxf3hR4VfqL0p+V/q61Lqj1D63fjuqNsVag8Whg3hVLkTUZDXmPpFEh0WhlBmqh4QLKGmECXRROCkzI0hSbgWRfDDTHio48xz97EDgZDjS3dw3RXLzExEWzogMlU6oMPyHpYpXLgqJD4egwiPoOhlQ2kx6DeTLDBmwqSuMXizQxGW+QWuHj4dD33RAFurZXkOczKMHbKJQwlJLB7AxKcJJGCbQug+kQJdSIDpXGSLPD92Uo2fJHK87jgUbv4tiKDgQlODt/pSs10FNDCbDU5mSDNMLZO/8L9t+FHSOXFJ2+vBW94lq98uTNg2AOSJM13D2+jKSI/4G/4yjhtMaME3pPYpCSN9+57DIO2Sko5s9xgRUSXlL651yz6WSnFs/htJYGd+d2sBtgxnsgdzZCypHh1eAJyVDYHkO5I5wUyHEZxqjh4/f1HwkMSdFrD3Cukcbw9WNPObdmn1iNrQHT8fA3AAD//zTwK9PfBgAA")
	data["core/01_basics.scm"] = decompress("H4sIAAAAAAAA/7RW3XKzNhC95yk2yUVECv6cXsbTGb+Hx/NFBuFoKiQiyZn67TurFSCw7F6Vm8Ts2d9ztGK32+2AKwGNseIDTtzJxhUFa0UntXiDt8G4WuruDdgveN9sYbvZlmVi1+I82esZMCO8vQh48YlPx5UT8NIlr7RU8MrQb7cDd9We/1N/X4yX+gxWfF+kFQ6knx0aoxvunwoApnh/ajk0RilXAOAb4WsrGjgQqpZaCxtM8zO6seAHX4K35QoCwGQHTLpa9IO/UoZbUADyYVBXUNL5O6HGwoAdOmCdtM7HgMcslp6DBWbFBIVjPvC60u4eLEDTqYB9UG8Oj79cLKorwT7yXMXhwyB0G/KNE+hKesYZLHMRNYcjIhJFKX4SCpnveWNNpDyyqXkvoDO2H+v6nNVQBWOF1iOEH6vA0tXYXW1sjUwm4mI/XF0ExRxHHQaRvMcHxV7MnLgQJ4IoF0rg0PO/RT1Xn+RB0YtkqCxFpUjswi1mP6ur4ZbsWW2hppp2QqwJTNsbJ5EXVCw25itHaYTYqIxHPh3aiXkc88gA/a2pxyIO4N6sAp0n017nRJ9TpMBvsRgPyobeTzGqffAviZzbQhKF/T91hIDZchJlUlrgzgnra27PSPyBlbg2j+HfRvGLG3XILJdOwHOCBwI46C/Ow0nAwKUV7XMQyCEai1i67KBio8M2ZZGlIas9bYG4oAhfLliPlczR3pHxm8Zi++GUhrFsojiLWY7KNFwFa3Kys4ySa4oiA53+0XxbxJP4ppMKG2hMP4T0n0wbD0x8Q0W2ak+2DD26jbTgFljwEqeTG3bYB1y3L4uZH4v0OKJ5edS5bqf5T5NfQNAnO+sn9N1MXnOPFDPhceVp7B3V/Vd3xj5qzthlb+vfzNjbVvOdGZttLAS43xdrRaO4RdpRIqS66dKIxtWmTU7+zRfB8irOQvB5Zdn9SDs0LN1XrsSvmP9tvDEp3v1bd+qGJhbh0xX7yU7iLDVOdAROkMxg8BdspuVGfNJBCFs/GCgyW66nkbr57K6uSri9ZeI9OXo7b+FZ6h+uZBsX8gc8U87lVZxNkOzgfX471zdrdxoBMfgj9ByV/QWsN23cEX+WsF3jTds+gL+v4eHLeMbLLmLD+/CRvPagL+eMCxnQcemimwT9RwQvC2lFiqmzGG1qM8AGfi+jt0J76a/RfYwyA5Bcz7VX1yRDnPnvie216j7Ac3sWHnrhvwzuKpQUnWeG76torwhQQrUPiLL4NwAA//8DaLIgzgwAAA==")
	data["core/02_binding.scm"] = decompress("H4sIAAAAAAAA/4RS0W6DMAx85yusPpkJpm6PRZP2HwhpAdwuWghTklbr309N4kBQ1/mlFT6ffXdpmqYBoQiG2dABeqlHqU9FgSMdpSZAaev4sR6UOFuC8FMWACj06BEXGtxsVh1f+AavgIr0yX1yb2lKW6t5EAow8u7LsixXmyfxRbzbwkWoc1iqyNWGBmiX0yyzBmolpn4UgGkm687mnqq7WI9nkUpa9zcs46bp213/wWbcm2PwKE1a9pAheZmMQkPLqK+uCLusJeNqYZJd2VyYALTOwE7qi1By5PdwgF12DMrjOveNULznFB+UAq4nMZjZx/kE6Ypn6OfxyjlD22/fAf9hVTcUtvErG9fHbqhW04/zYG9NX0LHd8VxLSYCQGbZZ+MArT9+6b+U3Upt9CJkfluVB/YRhFR+R+WZOqjeWeZj4Db64FblBUWOUMVvAAAA//9KoHJsyQMAAA==")
	data["core/03_branching.scm"] = decompress("H4sIAAAAAAAA/5xWy27rOAzd5ys42VQu7GJmlg0G6H8EBarYdCJAllNJySR/f0FKVi3FSXuvFgUqnmMeUnxks9lsQGqEdrT4CjsrTXtQZr9aiQ57ZRDM6EG5ppfaYZWum0G2doST0ejcCmArPDpfwdcxSr8nA/gDmmj9EKqHmm+N0lCzqYACkrc5tKabCXyj4/8DmgUVuYZ+tENVKKjpknBzBS8MdZl/scO9MlC/RRNTFnU0ZvR/pIWzwYZvtBCw0LOgpR1Nx9+pku9Wy5NDei8QGv10AaJX1nmIZnZPEOkcWt9Iu3erGIWQpgOhXHPG1o82UqrV7OFB/Af/0vfN3h8mQJUgwnkLa2XOUquORUbMK6wTOgogjWLLMQOIKPbvSV8421CzyfwPpQLKUssE1oGSi2Yl9ZuwOMsEHcrsZgOtdAgWP0/KogPlU4OIQR7/gv5kWmhHrcm5CH3BicLh6K/JMnfnQDBryj5Bqvnnoha+n5OfBP1fPLdggXg5WnhhsY6VaPSNxRbE9iz15HyPxl0HWJ+lXs97Np0sw7A9WuwaragNtBx2nQShqSJHSxljwQQB0rlAjV4n6oWo+Ak1KaovVfXlNedSFM/Z1fSRVaGXivxJWKkcwtqMIX5ox5PuYIcwSN8esFsX4iBvieL8oENKxkLD3EC+75+S8rCdSnDeXVQS97prKdyFZoObcsjPbf99yzB4YR9Zq73fCyk0Mf2hvHEV8mR9lLMQ0VfZ/gYhYB+Cl+ZHCRHyeNRXfoFnoIDDKIlRfoTiCh1APZvylRNDHy90u+objWHN7JTplNmHFbtKKSM71JORF+dsz81JYdn+fOZH7oOhPyHuTf2oLsKoNBMjm/w0pCBpLQd/3MvRms/9kF7i8/jPiFxPwUTR8y+LammD8jYv08wL+m6aZ+t7unuB3dhd71LSImcUqfgVAAD//9b094CUCQAA")
	data["core/04_predicates.scm"] = decompress("H4sIAAAAAAAA/7SV0a6bPAzH73kKfz43QSqfdr1q2oMcTToBDMoWEpqEbuzpJyeUsp56VxsXtiE/B/+Nm57P5zNoS9D5QB9hDtSbTieKVaUspSZQB+qVHzd6nu1awf1SVk9tr0ENi+tAhzHWx2UAZQZQJjY0zWl9BvCVwkLvn6rFWYpx21sNJsRUdqhh0DY+SdlT79VCyQ60J9d1/aU65L72NBhHzeyjLM3piR6lWUrwyqsNr4KK6wQqppBhwM/IL3qo7E2Vt8HpkLi/6QX+h/DypEOgTgdJp6Kp8y7C8MIZsipH419Shf/hP1e2J1y1XagG5XwCVfSWR/WzPX7f8F1fam7MVl4z6S74/W6fdnhoyJtqaTQOTuo+Ho/M7doZR+OB4VqrW1f4R+AWazcZXBBdyg04Yx/RnxT8Af20kR8O3KF2E7l/5HoKgCXCWuTsmiG7ikzyEyBbiWi9t6Qd4BZIXP4UyFYmFpeoZygHElcOEMxOZK7kANlKRD43ALOTGON6+sEVbYHEfaP1uw894BZInDUxAbIVCd9pC5idxJSxxexkZp658OIlyvF3c/I3czQ2xg2AbEWIRxnZysTU8jgWL1G+/UpdAixepPoe0PeipFmbAMhWJHwsov6gaQ5+MjweWyBxl0VbMxju9B5KbKDo7ZXRWySTVwrRtJaYvcUSHekCGOkirs/UGR6qLRC5FIwbAYsXqXVqPW+WvUTx3zggW4m4Upd8ACxeovLZh2yxrn4FAAD//1RDHNieCAAA")
	data["core/05_concurrency.scm"] = decompress("H4sIAAAAAAAA/3xTwY6bMBC98xVPysVelUq99JAc2v9ASDVmSJGMJzWQbvbrV9iYtQkKl9gz8/zevJlcLpcLlCFodnSGZqtn58jqR1GIlrreUjko7RjiyviOhtuHLIA/4spvEEYNTasgJL799ikpn3DdPM2OMqyhCdXN8dCPdIJYT0fP1QWWbyGPZScpQzC57zlbMuqRUb7g+IJDdOw04a7MTAtO9N3G++sr7CVlt3B5FmLUR67DK/Og4MKSKR2Ns5lOEA1dexuF1ZFp1Vw60nA0srkTRMBEfuy1PqV9yR69y++jqTh5NNsrWXJqepruG0Sl/yp7AsTyKyNDnRBW2rCf/9kf4AEyq6ChnxbdZ38IBajlthSxNnj52salLDAmXW/NrU+eR/q38qTtrksz3tR/WwCV6GarI8ZHsUTw46enC3kMDb+XY/9BB5VbDpZLvh3CMLDtJ3YRnhhryQRnE5I6G2Y1qN40/J70ZMnIfARANZJtyWUe+7p6cynx2a+Q5ju59I+0yg58+50KX+zkMBkkyFoWnwEAAP//bjriE48EAAA=")
	data["core/06_sequences.scm"] = decompress("H4sIAAAAAAAA/6RWTVPjMAy951eoveDsThnOdAZ+SKcH11HAjHGK7XSXf78T+SO2k1Jm6WU31pP03pOVsN/v98AVghgMPoJ1XHfcdGDxY0Qt0DYN67CXGoFZ/IALVyO2DQCTPTCuOzp+ns/Tj23w/ew+YyjG6Cn8X0vVtmWDTdVhWZyOy9r+HyoXQYZLO1V0BrY+Kri+c8CFA26BJ4GPsI0May67p+H0hsLBPYhBKUuk+PmsPiEEwpMYtOBuE1DLMkratSJ0/M0SFxRuMCtFQuB2GX4+o+5+UQjuvWZfiCn+fur4zqAAj9pJrdEAI+yMTPanm2DzkbMy2T/5fqyXxrqY0wIzmD1mNSa0P8h75zrEoO0tFRPmPzWUqdpW1D3BH/J/W6FP60SOPfsizWyq+oQ4Pc+prF0AvT0rsIKEQv3iXjepU+GeDxb+nQ1eCvfEMGqHFVeK/iZwbJEZApX5dSYFy97kcnBcauFp1Dn+cPb8IVeq3eumobqOtB2s4LopW5LyKHSw8C6tlfplwe0r4gnwRCUe1uKEmfqXwjoUU057rW9K9RdxYWiKx/y2PTaZMtIdfKDfIUldzEXqDv+ujZUQOsyzTp2FJQ+Tq218F2+pNgyjg6GH0zDqzm7bRLbgBR32fFTuB/zWSnzNMya07bFaFh5tDzcJDra838HvP6+opxtQzX/mRbMuN6NeTG7dtbXk1n29lPQ5vrqZAVc4WNC8rqyYwULfd1Qu4ItVvvVecMAOfbEEx9VFOZTvjeMK08xJA/3NF8ody2e0CzfGYDeK6S+ZA+tHLUBq6XLuxew8OE6vht6Sz8r0uV3xSshlVwUmcDqbux+bmX7O/AqVQAMIn7e+0pkKT9uU3W+DFzQWsys+NfOn8qSw+vgFeH1INOhDesfadMX+BQAA//9/6ihdzAoAAA==")
	data["core/07_lazy-seq.scm"] = decompress("H4sIAAAAAAAA/6RXy27bOhPe+ykm2fxkf6uwt02R9D0EA6GlkS2UohSSbuM+/QEv4kWifAocbRJr7t98M6ReXl5egHGEZpT4DTj7cweFHzcUDardjrTY9QKrgTVyBGLElcIP+Arnsb3THcB7ePnFyIdzy4BQ2P+wCpQGH0A0+2ni3ISGZuTcWBNvUklswMirXgiUQHI184Q4/jcA6TsgTLRAXr3bAw3C8JAnHCZ9f3POcgXSjEIB6XqpfDCX5ZxFi43zTIFInHUWTv5HqH+VJL2su/p97TnCJLGNxWcVkd9XFBVHDbUCYkC2eqdQvpF0IVsVJN4UiHXeqSUGrkhjuUzE1aSof5KUWzlOy1bl2aaNM9rbjfON2ugRSY038M4M0lcp4DH9ytNwYlL3uh/FDqBeJ0aCPPgJGqddYqM0Tpnhkoke/9CzXQF/I61eea/0ehAKrC2kZ9NwnYkZmee0Ll0ycUFbwpyzfQMHED2H41wfZ0ovFcy7uWP2/wOFI1RHOhs5/qWmibIVJvV4t9EGjmtptCzHsfUG8M0cNMO0W0DcQk3EjfM3l5qDXTOh+R20vGEyL+GpyXeH5YECwOtSXLL4hlzh/Ov7yiJGsag0wwQrvFJaOGFcrB4R8n8vsZUnIORcqU/F7g9scgjeRJMRNx3agU2V6sWFo0klY+161f4rySHZpiZqulOL9LYGaQ7pfo1lptMYyoGv9o/aKmtiknGO3BWm/rIyNk2eKm82M0gPDlWqwnKR1J0f7l/Y6FE621i+oiXm+f557VC7oguNU3k3OKRdxhaWjjo0Q+kyg9EZWOzmnNJ13/Vco4SMLw/2vVOPGz/b9Y9pEiFLCFLGx4ATSVFCwo6YL7/ows9YXqDcICTJypIJfKsDxl+HulFWyJqruTFV+DlJtbgYKf1ktWAfNUo3o2YUDdMprTO8nTjFWz24FeUNWF8HCh1Qf9MCtd0DE+xhC3ysDZ2omZVp7Gz0jtp2bPqPVgnr56mNKP/ppwxiNz5m/uyxvDYY2GTaYlmzNvRNS9yE+VEbbCkSxZ1oZxPuJ1bnXrS9uKio6ntDmFIodcXkRc3NnzfUmcIz08DRnBOjQDCxvCfoFUj8uPUS22eHoYn4BUitr31ozTmnQC3wU3v4zzQ/CWt1HwCItT4sZSZlLzsmhLE08dka12kz3+0ODB8Pe3UfwgeE5dYT7I1bujKy+Od2btxs9vOk5T5oqTcNkxpVz0Q1ybG9NfkserjUfaguKCCG/KRALigMHkRpCc/N9AyflMYtnoBTc9SVug/+Im7ciou+zrNlQZhDuGsAzY4Oa/+LceVV0yPHUDu49wwsxjekSO+is8/luVPK36y7/5K/sX+Qf3A/F2JKWNov8599bubvOPNuh2wfEbAsiQ6tr/2POQVHkn8CAAD//1YfRhgcDwAA")
	data["core/08_threading.scm"] = decompress("H4sIAAAAAAAA/+xVy27bMBC8+yumzqHLQiySHmvU9n8IBsLYa0coRTkk5bZ/X/ClyI7ipodegvIigNzHzOyQWiwWCyjN2HaWv8I/Wla7xhxmM9rxvjEMSnvS8ZNc6sZ5eGUP7MUMoN5odg4U9lejg7LoIqGsvCOGNrJVW9tBLmdATSelexaIn83zDj5j39nW5Uqk2X8C1XoCI+0b63wOF5tRa6Del2MtML0uEizI8rV4bAZ29ySXoGqPKmGu1lagWqcCGY/YTDB/H9Qz97XN/N/C3XUtv4vJa/aoT0rfZPLjYvTjkY00nQeZXutV4HgjZlMlKQkSlUzVpjz0upL/pbyU8tmVMfRNWgYIcomHxoQHcSwbKefYeqnswcWu1DiZ4+RWq95xSROYlwL5oO2dxwND4agayzuceOs7O0+FPnB79L9WudcZlbny0KycR2c4v9S8i5FoHCw/9aHePCIM8lHtPB/xF4OtDf/0wFib64OtjWoZoMLxdnq2o4TkwSHh7pUEkXXdg7IiAdrzjLNBqtg/OwRV4Hs9ZmySPOASEJI3qNaxUVijf2CrvrPcdmZXxpvlt0ls1T7sFCidZeiXHgEoTXqVnZCZz1/4YmwIgL7hS4BqDv6xZIoXWclL89w7Tf8Y3JVB4fbs/tbRNeXsTgw3754GNielRZxAlQpVhXO4Q6hChfi6iHOx8u2JYv3pIUr9y1MU/5xpTiOo1ZpadZwYwUe5FEOFifubEPxrCOcYfgcAAP//ImobVk8JAAA=")
	data["core/09_functions.scm"] = decompress("H4sIAAAAAAAA/4RUTW/bMAy951cQuZQqku6+AGt/wLDDdjQMVLXpQINNJZIMzP9+kGTJsuelOiXG48d7fOTlcrmA7AkabegrdCM3Tmm2hwO21Cmm8yAbowF7ch3Dh+JW8dXCC3zodhIHAMReDh+tPBtq4CaNpXNGoR4dKPYw/1B1gJbuxScImZ8Bq46/K+sAO2VsCILl1QdYv6rjX9MACR1jxQP0DzlQQqOhJUbsoX9Shq3BUIsCjtJaMu4szdWusqDkFrBX1r2m0G0VANQGkO4QqTwtKu5gA37Byp6+FPi95L1uZO+re+YbxDEPaBitg0azk4qB5UDtYoHjiup2svJ2I27BDzjJG5nWYtZMsSjqvvs5B4+cvAXO37w6PlzA6S14KYKrOptMiMWFka33l1OyPwBU6BsVod06/4cXCDMOI5mL703Jt99PrxBzHOe0YOg+KkMWZJYhqTC7HDy1nVQBn35HbZ4BG822aGgJFkKIWmyXrNHDLVARoFpip9xUP6Q6cpNZzmsU6qTWrsR2GuDoPx7zQq0dXylmMtk2fipPBaXQ7SYiMtKjI5M3MLSyYDYRgXgMSBu1CSg2a741VXlYUrHUUpSyFGBt/3RpHsMCdMXjc/xKqUWiT9CBc5H8v91s0q7u03sy4emfi3OaaZSjiYNdpN83Hf4e/7jSTF5+iD7aGCi2st9FdXrDQd7ynmAnNm6aCc2Far/vfwMAAP//h6Qqyn4GAAA=")
	data["core/10_exceptions.scm"] = decompress("H4sIAAAAAAAA/6xWzXLjNgy++ymwziFUJ2qbHpNpk/fgeCa0BDnq0pQCytl1O333Dv8kiyJl77Q6JLYFgB8+AB/x/Pz8DEIiVB3hE+D3Cvuh7ZTebJjEoVHAmRTHfS1KwgpaXVZCSmD6fIRKipPGYgPTw4SqgcmuEvIF9Pk4e+ksZKuHl5Sve40fxg9Y05IegllRFJvNFCHCM1Tv5b5VdasOwJqOjglIn1gNHb3A4rWz+R1+AyZRHYZ3Z5IC7rKyJ8DjdUzAHHroBWmsE6ACm/fOPsvJMks1BHODZOnxRePHC7CnveyqrwHAGuKmVULK8w9gDh5Z1AGDSy5g+B+w4veeVoF+6egHyU0klAMw0LkUfY+GDHc4fMXzt47q5EBUnfrTgwQeDN2vLHz1OQT/3drRPeF/P1uH5okhXDnbGAXm9fysk5KoNTBbzJRFeP72hb5nBfjesB8D/Xz3T8rNqNFPwHgzlwZd7FLW/uEEjPDCGFate2BTmlQA7JIZOGbrXCQ+zis00BfG/LJwfci6KWZuWWQ2XmDHRmQXPdhPzM0DrsezM+ThLfC5At2M7wmlxhk9JFqNsD0KacQSa9s9TsI82K3ptLVec9Y9Yd1WYkBg5iMgURlfK6EzJLDhnVDUpcaP8g9zz4BxSnbI1EdynmYyz7GPYuNFh7jpavx/Dxfollz3JNR0a+h0rkJrpKEUdNDRwbPBg63rMR9TdQP0oiWstxF14q+zoSuK5RkNIjEfuCxLAPxTkA3AxssJfi3yvAJw24jg6b0cVqNFsYtj1a4QcC8k/iJxWBloh8fzuDb5wBshrbSZ+DbyHg+tAoMugcOjmWoXahapTSjhLeUfQ6zWP1Owb++obAtkpXe9pjlyuJ26eUEfi6QucpYa2+zULjn0vRqnf0Un9l19Tl9Jb94UWDHW11X14XXM+2pdMtedxAE4EgE7oDIjvkWibczjBOEBiRLLR+YeeXhlou/lGWynJ/rsgqTESjVqMh/ohGDO3t0kt+MWmNqqTMqMu7sB4n0tOd+csOo+keIlMG1cSRTq1AMbLzRvntLYBWs8CKANklrfLXzb9ndBbsIQOJ9Y2/N3Hjfud/GGNeIOKuDC7tIC8sZqbJAy0jK17gNb1scenyx85Pvq8i2KYrfJMOZrtAz25ikjrCxhrpQOjp062wAzgGOw/FZmhsaFuxuvhutPPhyhDuEebw+XpI61DXh0dn8xoQv3N0+fJ+FtlBb3Q8J+tiapVu5CVUwjtArLo6ios9MOP891PNECF5tqMDUD/m8AAAD//wCF663QDwAA")
	data["core/11_io.scm"] = decompress("H4sIAAAAAAAA/9xTzW7bPBC86ynm03chBStNe7TROu+hCChNrW2iFKmQdAy/fUHq33EL5FodbGjJmZ2dHe12ux2EJkjraAv1xWYZa0hq4QiFMgUKewkFCnKu4OnsqAyh8J2QVADIX33O57qhq1aGCuSvJl8AWOfKVnTlVYVzaS5a43gxEp7eeAawVnRgWrSHRoC9C30hDqaOYPHqHkMl/YElZF/iEZ2eyLTuhyccrWt9aqApoGpF11HzQIsP7r/hcp0I2fVMBszTG3rQ1CcebgdXtlenQhSknA/jxUHS3xmO1pUk5BkVaWrBHM0E9eLin5uOG4h4Hp/V7GY1vOg6fUPnFqV7tnFvdzzKhM/Z+E+5qEzQj52MvnzWTKEcNeU7yWDdGOnEaRqwD+VebgTt56zPbMl025EppbaeFnRpRRJsuzwY9hG/qDTCHpJDwtjSdkvishXS2QU9Dso0ypw8nnCwza1X7D25UAp38tkkczHbiOHIHxC1Fx9wIAj8olvZf9M9Mk/s0pom0VbsexzHnMJ5Qfk8uPOTHeikDDYvSVedDZgfj0DfJlC0h1UbNsl55li8feUfQ7OOUpV8/X9K1P0mVtS8ngPLgrst0ztbU21e+uz2v5NsXo/TrVKvjND6BtYLiamtOc9+BwAA//8zmRAOywUAAA==")
	data["core/12_os.scm"] = decompress("H4sIAAAAAAAA/3zPwUrEMBAG4Hue4qe9zASr9Wo9+B7LgrGdLoF0KpPU55c2i0UQcwqZ7//DDMMwICTBuJq8YM3O0SRjCibwol8ePtgtez7e56jSLWG0FVTiInjEvNqS2QHvlKR40CWXYKUFaNzMREu3S8bV4edcTPKWSgv6kFtUPLzVml9GdGr3y38902bVdDh4/ZuvXBHFGfSKQz33fc9nlj4taklah43m5q8hPZ3hPY9myQ3f5X0JZvcdAAD//+SGnjxGAQAA")
	data["core/13_alist.scm"] = decompress("H4sIAAAAAAAA/2xQS26FMAzcc4pZOot3AqReBLEwwahpTQJJqorbV4QUKGpWkWc8H7ctWAU2RMELrC7lpqFRJucFxCkFi0/ZYIOqaQD6fhcPSrKeI4BUMjpRmUGTiykfWF9AgL68Skog9uOxulNNRf88krXYVZnCM0/mLRZF+XV7suquudfZ673ewvAhNl+VeFl0wzG9RGjmxXIGKc/DyKAiiO4erfqXf/9PIZqcZok4z7Xn+QkAAP//QzCIXnMBAAA=")
}