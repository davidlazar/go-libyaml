package yaml

// #cgo LDFLAGS: -lyaml
// #include <yaml.h>
import "C"
import "errors"
import "unsafe"

type scalarNode struct {
	value  *C.yaml_char_t
	length C.size_t
	style  C.yaml_scalar_style_t
}

type sequenceNode struct {
	items struct {
		// actually *C.yaml_node_item_t
		start *C.int
		end   *C.int
		top   *C.int
	}
	style C.yaml_sequence_style_t
}

type mappingNode struct {
	pairs struct {
		start *C.yaml_node_pair_t
		end   *C.yaml_node_pair_t
		top   *C.yaml_node_pair_t
	}
	style C.yaml_mapping_style_t
}

func Load(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var parser C.yaml_parser_t
	var document C.yaml_document_t

	C.yaml_parser_initialize(&parser)
	C.yaml_parser_set_input_string(&parser, (*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	if ok := C.yaml_parser_load(&parser, &document); ok != 1 {
		return nil, errors.New("failed to parse yaml")
	}

	root := C.yaml_document_get_root_node(&document)
	return LoadNode(&document, root), nil
}

func LoadNode(document *C.yaml_document_t, node *C.yaml_node_t) interface{} {
	if node == nil {
		return nil
	}

	switch node._type {
	case C.YAML_NO_NODE:
		return nil

	case C.YAML_SCALAR_NODE:
		ystr := (*scalarNode)(unsafe.Pointer(&node.data))
		rstr := C.GoBytes(unsafe.Pointer(ystr.value), C.int(ystr.length))
		//tag := C.GoString((*C.char)(unsafe.Pointer(node.tag)))
		return string(rstr)

	case C.YAML_SEQUENCE_NODE:
		yseq := (*sequenceNode)(unsafe.Pointer(&node.data))
		rseq := make([]interface{}, 0)

		start := uintptr(unsafe.Pointer(yseq.items.start))
		end := uintptr(unsafe.Pointer(yseq.items.top))
		off := unsafe.Sizeof(*yseq.items.start)

		for i := start; i < end; i += off {
			item := (*C.int)(unsafe.Pointer(i))
			itemNode := C.yaml_document_get_node(document, *item)
			elem := LoadNode(document, itemNode)
			rseq = append(rseq, elem)
		}

		return rseq

	case C.YAML_MAPPING_NODE:
		ymap := (*mappingNode)(unsafe.Pointer(&node.data))
		rmap := make(map[string]interface{}, 0)

		start := uintptr(unsafe.Pointer(ymap.pairs.start))
		end := uintptr(unsafe.Pointer(ymap.pairs.top))
		off := unsafe.Sizeof(*ymap.pairs.start)

		for i := start; i < end; i += off {
			pair := (*C.yaml_node_pair_t)(unsafe.Pointer(i))

			keyNode := C.yaml_document_get_node(document, pair.key)
			valNode := C.yaml_document_get_node(document, pair.value)
			key := LoadNode(document, keyNode)
			val := LoadNode(document, valNode)

			rmap[key.(string)] = val
		}

		return rmap

	default:
		return nil
	}
}
