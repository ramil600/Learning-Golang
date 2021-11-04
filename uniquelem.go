
//SingleElem returns index of the only non-repeating item  in a slice 
func UniqueElem(elems []int) int {
                insertedsingle := make(map[int]bool)
                
                for _, v := range elems {

                                _, ok := insertedsingle[v]
                                if ok == false {
                                                insertedsingle[v] = true
                                } else {
                                                insertedsingle[v] = false
                                }
                }
                for i, v := range elems {
                                if  insertedsingle[v] == true {
                                                return i
                                }
                }
                return 0 
}
