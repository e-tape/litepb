package t_pool

/*
#include <pthread.h>
*/
import "C"

func getThreadId() uint64 {
	return uint64(C.pthread_self())
}
