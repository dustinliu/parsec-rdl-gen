package utils

import (
	"fmt"
	"github.com/ardielle/ardielle-go/rdl"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
)


func JavaGenerationHeader(banner string) string {
	return fmt.Sprintf("%s// Please DO NOT edit directly; changes could be overwritten.\n//", javaGenerationBanner(banner))
}

func JavaGenerationOrigHeader(banner string) string {
	return fmt.Sprintf("%s// WILL NOT be auto-generated if file has already existed.\n//", javaGenerationBanner(banner));
}

func javaGenerationBanner(banner string) string {
	return fmt.Sprintf("//\n// This file is generated by %s\n", banner)
}

func JavaGenerationPackage(schema *rdl.Schema, namespace string) string {
	return JavaGenerationOrigPackage(schema, namespace) + ".parsec_generated"
}

func JavaGenerationOrigPackage(schema *rdl.Schema, namespace string) string {
	if namespace != "" {
		return namespace
	}
	return string(schema.Namespace)
}

func camelSnakeToKebab(name string) string {
	s := strings.Replace(name, "_", "-", -1)
	result := make([]rune, 0)
	wasLower := false
	for _, c := range s {
		if unicode.IsUpper(c) {
			if wasLower {
				result = append(result, '-')
			}
			result = append(result, unicode.ToLower(c))
			wasLower = false
		} else {
			result = append(result, c)
			wasLower = true
		}
	}
	return string(result)
}

func JavaGenerationRootPath(schema *rdl.Schema, basePath string) string {
	if basePath != "" {
		if schema.Version != nil {
			if basePath != "/" {
				return fmt.Sprintf("%s/v%d", basePath, *schema.Version)
			}
			return fmt.Sprintf("/v%d", *schema.Version)
		}
		return basePath
	} else if schema.Name != "" {
		n := camelSnakeToKebab(string(schema.Name))
		if schema.Version != nil {
			return fmt.Sprintf("/%s/v%d", n, *schema.Version)
		} else {
			return fmt.Sprintf("/%s", n)
		}
	}

	return "/"
}

func JavaGenerationDir(outdir string, schema *rdl.Schema, namespace string) (string, error) {
        return _javaGenerationDir(outdir, schema, "./target/generated-sources/java", JavaGenerationPackage(schema, namespace))
}

func JavaGenerationSourceDir(schema *rdl.Schema, namespace string) (string, error) {
        return _javaGenerationDir("", schema, "./src/main/java", JavaGenerationOrigPackage(schema, namespace))
}

func _javaGenerationDir(outdir string, schema *rdl.Schema, defaultDir string, pack string) (string, error) {
	dir := outdir
	if dir == "" {
		dir = defaultDir
	}
	//pack := javaGenerationPackage(schema)
	if pack != "" {
		dir += "/" + strings.Replace(pack, ".", "/", -1)
	}
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0755)
	}
	return dir, err
}

func JavaGenerateResourceError(schema *rdl.Schema, writer io.Writer, namespace string) error {
	return _javaGenerateTemplate(schema, writer, javaResourceErrorTemplate, namespace)
}

func JavaGenerateParsecResourceError(schema *rdl.Schema, writer io.Writer, namespace string) error {
	return _javaGenerateTemplate(schema, writer, javaParsecResourceErrorTemplate, namespace)
}

func JavaGenerateParsecErrorBody(schema *rdl.Schema, writer io.Writer, namespace string) error {
	return _javaGenerateTemplate(schema, writer, javaParsecErrorBodyTemplate, namespace)
}

func JavaGenerateParsecErrorDetail(schema *rdl.Schema, writer io.Writer, namespace string) error {
	return _javaGenerateTemplate(schema, writer, javaParsecErrorDetailTemplate, namespace)
}

func _javaGenerateTemplate(schema *rdl.Schema, writer io.Writer, content string, namespace string) error {
	funcMap := template.FuncMap{
		"package": func() string {
			s := JavaGenerationPackage(schema, namespace)
			if s == "" {
				return s
			}
			return "package " + s + ";\n"
		},
	}
	t := template.Must(template.New("util").Funcs(funcMap).Parse(content))
	return t.Execute(writer, schema)
}

const javaResourceErrorTemplate = `{{package}}
public class ResourceError {

    public int code;
    public String message;

    public ResourceError code(int code) {
        this.code = code;
        return this;
    }

    public ResourceError message(String message) {
        this.message = message;
        return this;
    }

    public String toString() {
        return "{code: " + code + ", message: \"" + message + "\"}";
    }

}
`

const javaParsecResourceErrorTemplate = `{{package}}
public final class ParsecResourceError implements java.io.Serializable {

    private ParsecErrorBody error;

    public ParsecResourceError() { }

    public ParsecErrorBody getError() { return error; }

    public ParsecResourceError setError(ParsecErrorBody error) { this.error = error; return this; }
}
`

const javaParsecErrorBodyTemplate = `{{package}}
import java.util.List;

public final class ParsecErrorBody implements java.io.Serializable {

    private Integer code;
    private String message;

    private List<ParsecErrorDetail> detail;

    public ParsecErrorBody() { }

    public Integer getCode() { return code; }

    public String getMessage() { return message; }

    public List<ParsecErrorDetail> getDetail() { return detail; }

    public ParsecErrorBody setCode(Integer code) { this.code = code; return this; }

    public ParsecErrorBody setMessage(String message) { this.message = message; return this; }

    public ParsecErrorBody setDetail(List<ParsecErrorDetail> detail) { this.detail = detail; return this; }
}
`

const javaParsecErrorDetailTemplate = `{{package}}
public final class ParsecErrorDetail implements java.io.Serializable {

    private String message;

    private String invalidValue;

    public ParsecErrorDetail() { }

    public String getMessage() { return message; }

    public String getInvalidValue() { return invalidValue; }

    public ParsecErrorDetail setMessage(String message) { this.message = message; return this; }

    public ParsecErrorDetail setInvalidValue(String invalidValue) { this.invalidValue = invalidValue; return this; }
}
`

func JavaGenerateResourceException(schema *rdl.Schema, writer io.Writer, namespace string) error {
	return _javaGenerateTemplate(schema, writer, javaResourceExceptionTemplate, namespace)
}

const javaResourceExceptionTemplate = `{{package}}
public class ResourceException extends RuntimeException {
    public final static int OK = 200;
    public final static int CREATED = 201;
    public final static int ACCEPTED = 202;
    public final static int NO_CONTENT = 204;
    public final static int MOVED_PERMANENTLY = 301;
    public final static int FOUND = 302;
    public final static int SEE_OTHER = 303;
    public final static int NOT_MODIFIED = 304;
    public final static int TEMPORARY_REDIRECT = 307;
    public final static int BAD_REQUEST = 400;
    public final static int UNAUTHORIZED = 401;
    public final static int FORBIDDEN = 403;
    public final static int NOT_FOUND = 404;
    public final static int CONFLICT = 409;
    public final static int GONE = 410;
    public final static int PRECONDITION_FAILED = 412;
    public final static int REQUEST_ENTITY_TOO_LARGE = 413;
    public final static int UNSUPPORTED_MEDIA_TYPE = 415;
    public final static int MISDIRECTED_REQUEST = 421;
    public final static int PRECONDITION_REQUIRED = 428;
    public final static int TOO_MANY_REQUESTS = 429;

    // Those are non-standard HTTP code only used in Sherpa
    public final static int SHERPA_V1_TABLE_NOT_FOUND = 420;
    public final static int SHERPA_V1_INVALID_DATA = 425;
    public final static int SHERPA_V1_TEST_AND_SET_FAILED = 427;
    public final static int SHERPA_V1_TOO_MANY_REQUESTS= 432;
    // End of non-standard HTTP code

    public final static int INTERNAL_SERVER_ERROR = 500;
    public final static int NOT_IMPLEMENTED = 501;

    public final static int SERVICE_UNAVAILABLE = 503;

    public static String codeToString(int code) {
        switch (code) {
        case OK: return "OK";
        case CREATED: return "Created";
        case ACCEPTED: return "Accepted";
        case NO_CONTENT: return "No Content";
        case MOVED_PERMANENTLY: return "Moved Permanently";
        case FOUND: return "Found";
        case SEE_OTHER: return "See Other";
        case NOT_MODIFIED: return "Not Modified";
        case TEMPORARY_REDIRECT: return "Temporary Redirect";
        case BAD_REQUEST: return "Bad Request";
        case UNAUTHORIZED: return "Unauthorized";
        case FORBIDDEN: return "Forbidden";
        case NOT_FOUND: return "Not Found";
        case CONFLICT: return "Conflict";
        case GONE: return "Gone";
        case PRECONDITION_FAILED: return "Precondition Failed";
        case UNSUPPORTED_MEDIA_TYPE: return "Unsupported Media Type";
        case INTERNAL_SERVER_ERROR: return "Internal Server Error";
        case NOT_IMPLEMENTED: return "Not Implemented";
        case MISDIRECTED_REQUEST : return "Misdirected Request";
        case PRECONDITION_REQUIRED: return "Precondition Required";
        case TOO_MANY_REQUESTS: return "Too Many Requests";
        case REQUEST_ENTITY_TOO_LARGE: return "Request Entity Too Large";
        case SHERPA_V1_TABLE_NOT_FOUND: return "Table Not Found";
        case SHERPA_V1_INVALID_DATA: return "Invalid Data";
        case SHERPA_V1_TEST_AND_SET_FAILED: return "Test And Set Failed";
        case SHERPA_V1_TOO_MANY_REQUESTS: return "Too Many Requests";
        default: return "" + code;
        }
    }

    int code;
    Object data;

    public ResourceException(int code) {
        this(code, new ResourceError().code(code).message(codeToString(code)));
    }

    public ResourceException(int code, Object data) {
        super("ResourceException (" + code + "): " + data);
        this.code = code;
        this.data = data;
    }

    public int getCode() {
        return code;
    }

    public Object getData() {
        return data;
    }

    public <T> T getData(Class<T> cl) {
        return cl.cast(data);
    }

}
`
func GetUserDefinedTypeAnnotations(userDefinedType rdl.TypeRef, schemaTypes []*rdl.Type) map[rdl.ExtendedAnnotation]string {
	for _, schemaType := range schemaTypes {
		switch schemaType.Variant {
		case rdl.TypeVariantStructTypeDef:
			if userDefinedType == schemaType.StructTypeDef.Type {
				return schemaType.StructTypeDef.Annotations
			}
		case rdl.TypeVariantStringTypeDef:
			if userDefinedType == schemaType.StringTypeDef.Type {
				return schemaType.StringTypeDef.Annotations
			}
		case rdl.TypeVariantMapTypeDef:
			if userDefinedType == schemaType.MapTypeDef.Type {
				return schemaType.MapTypeDef.Annotations
			}
		case rdl.TypeVariantArrayTypeDef:
			if userDefinedType == schemaType.ArrayTypeDef.Type {
				return schemaType.ArrayTypeDef.Annotations
			}
		case rdl.TypeVariantBytesTypeDef:
			if userDefinedType == schemaType.BytesTypeDef.Type {
				return schemaType.BytesTypeDef.Annotations
			}
		case rdl.TypeVariantNumberTypeDef:
			if userDefinedType == schemaType.NumberTypeDef.Type {
				return schemaType.NumberTypeDef.Annotations
			}
		case rdl.TypeVariantUnionTypeDef:
			if userDefinedType == schemaType.UnionTypeDef.Type {
				return schemaType.UnionTypeDef.Annotations
			}
		}
	}
	return make(map[rdl.ExtendedAnnotation]string, 0)
}

func JavaType(reg rdl.TypeRegistry, rdlType rdl.TypeRef, optional bool, items rdl.TypeRef, keys rdl.TypeRef) string {
	t := reg.FindType(rdlType)
	if t == nil || t.Variant == 0 {
		panic("Cannot find type '" + rdlType + "'")
	}
	bt := reg.BaseType(t)
	switch bt {
	case rdl.BaseTypeAny:
		return "Object"
	case rdl.BaseTypeString:
		return "String"
	case rdl.BaseTypeSymbol, rdl.BaseTypeTimestamp, rdl.BaseTypeUUID:
		return "String"
	case rdl.BaseTypeBool:
		if optional {
			return "Boolean"
		}
		return "boolean"
	case rdl.BaseTypeInt8:
		if optional {
			return "Byte"
		}
		return "byte"
	case rdl.BaseTypeInt16:
		if optional {
			return "Short"
		}
		return "short"
	case rdl.BaseTypeInt32:
		if optional {
			return "Integer"
		}
		return "int"
	case rdl.BaseTypeInt64:
		if optional {
			return "Long"
		}
		return "long"
	case rdl.BaseTypeFloat32:
		if optional {
			return "Float"
		}
		return "float"
	case rdl.BaseTypeFloat64:
		if optional {
			return "Double"
		}
		return "double"
	case rdl.BaseTypeArray:
		i := rdl.TypeRef("Any")
		switch t.Variant {
		case rdl.TypeVariantArrayTypeDef:
			i = t.ArrayTypeDef.Items
		default:
			if items != "" && items != "Any" {
				i = items
			}
		}
		gitems := JavaType(reg, i, false, "", "")
		//return gitems + "[]" //if arrays, not lists
		return "List<" + gitems + ">"
	case rdl.BaseTypeMap:
		k := rdl.TypeRef("Any")
		i := rdl.TypeRef("Any")
		switch t.Variant {
		case rdl.TypeVariantMapTypeDef:
			k = t.MapTypeDef.Keys
			i = t.MapTypeDef.Items
		default:
			if keys != "" && keys != "Any" {
				k = keys
			}
			if items != "" && keys != "Any" {
				i = items
			}
		}
		gkeys := JavaType(reg, k, false, "", "")
		gitems := JavaType(reg, i, false, "", "")
		return "Map<" + gkeys + "," + gitems + ">"
	case rdl.BaseTypeStruct:
		switch t.Variant {
		case rdl.TypeVariantStructTypeDef:
			if t.StructTypeDef.Name == "Struct" {
				return "Object"
			}
		}
		return string(rdlType)
	default:
		return string(rdlType)
	}
}