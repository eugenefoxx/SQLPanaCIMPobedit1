package gosaprfc

import (
	"fmt"
	"os"

	"github.com/sap/gorfc/gorfc"
)

func abapSystem() gorfc.ConnectionParameters {
	user := os.Getenv("user")
	passwd := os.Getenv("passwd")
	ashost := os.Getenv("ashost")
	sysnr := os.Getenv("sysnr")
	client := os.Getenv("client")
	lang := os.Getenv("lang")

	return gorfc.ConnectionParameters{
		"user":   user,   //"demo",
		"passwd": passwd, //"pass00word", //"welcome",
		"ashost": ashost, //"10.68.110.51",
		"sysnr":  sysnr,
		"client": client, //"620",
		"lang":   lang,
	}
}

func (g *goRFCStorage) GetSPP(woname string) (spp interface{}, err error) {
	c, err := gorfc.ConnectionFromParams(abapSystem())
	if err != nil {
		//fmt.Println(err)
		g.logger.Errorf(err.Error())
		return nil, err
	}
	g.logger.Infof("Connected: %v", c.Alive())

	orderSAP := "00000" + woname //"000001000880"

	info_order, err := c.Call("Z_IEXT_PRODORD_INFO", map[string]interface{}{
		"AUFNR": orderSAP,
		"UCODE": "21717",
		"PCODE": "NEWPASSWORD1",
	})
	if err != nil {

		fmt.Printf("ERROR Z_IEXT_PRODORD_INFO %v\n", err.Error())
		g.logger.Errorf("ERROR Z_IEXT_PRODORD_INFO %v", err.Error())
		return nil, err
	}
	//var spp interface{}
	for key, value := range info_order {
		//fmt.Println("[", key, "] has items:")
		//fmt.Printf("ref %v\n", reflect.TypeOf(value))
		if key == "PRODUCT" {
			//fmt.Printf("ttt %v\n", value)
			if value == "" {
				// log.Fatalf("Order is fail, not PRODUCT  %v\n", value)
				g.logger.Fatalf("Order is fail, not PRODUCT  %v", value)
			}

		}
		if key == "RESITEMS" {
			for _, v := range value.([]interface{}) {
				//fmt.Println("\t-->", k, ":", v)
				//fmt.Printf("v: %v\n", reflect.TypeOf(v))
				for kk, vv := range v.(map[string]interface{}) {
					//fmt.Println("\t-->", kk, ":", vv)
					if kk == "POSID" {
						fmt.Printf("SPP POSID: %v\n", vv)
						spp = vv
						break
					}
				}

			}
		}
	}

	fmt.Printf("dd spp: %v\n", spp)
	g.logger.Infof("SPP: %v", spp)
	if spp == nil {
		spp := ""
		return spp, nil
	}
	g.logger.Infof("SPP after nil: %v", spp)
	c.Close()

	return spp, nil
}
