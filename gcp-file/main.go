package main

import (
	"context"
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	//  GOOGLE_APPLICATION_CREDENTIALS="/Users/jay_liao/go/src/jay-test/in-sports-test-storage.json" ./jay-test -o=jay-test-fat-log:jay/jay.json  write

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	// read "o"

	var o string
	flag.StringVar(&o, "o", "", "source object; in the format of <bucket:object>")
	flag.Parse()
	log.Printf("o : %v", o)

	names := strings.Split(o, ":")

	if len(names) < 2 {
		usage("missing -o flag")
	}

	bucket, object := names[0], names[1]

	log.Printf("buckets : %v", bucket)
	log.Printf("object : %v", object)
	if len(os.Args) < 3 {
		usage("missing subcommand")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[2] {
	case "write":
		// ss, _ := ioutil.ReadFile("uu.txt")
		ss := os.Args[3]
		filePath, err := b64.StdEncoding.DecodeString(ss)
		if err != nil {
			log.Fatal(err)
		}
		if err := write(client, bucket, object, filePath); err != nil {
			log.Fatalf("Cannot write object: %v", err)
		} else {
			log.Print("ok")
		}
	case "read":
		data, err := read(client, bucket, object)
		if err != nil {
			log.Fatalf("Cannot read object: %v", err)
		}
		fmt.Printf("Object contents: %s\n", data)
	case "metadata":
		attrs, err := attrs(client, bucket, object)
		if err != nil {
			log.Fatalf("Cannot get object metadata: %v", err)
		}
		fmt.Printf("Object metadata: %v\n", attrs)
	case "makepublic":
		if err := makePublic(client, bucket, object); err != nil {
			log.Fatalf("Cannot to make object public: %v", err)
		}
	case "delete":
		if err := delete(client, bucket, object); err != nil {
			log.Fatalf("Cannot to delete object: %v", err)
		}
	}
}

func write(client *storage.Client, bucket, object string, filePath []byte) error {
	ctx := context.Background()

	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.Write(filePath)

	if err := wc.Close(); err != nil {
		return err
	}
	// [END upload_file]
	return nil
}

func list(w io.Writer, client *storage.Client, bucket string) error {
	ctx := context.Background()
	// [START storage_list_files]
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	// [END storage_list_files]
	return nil
}

func listByPrefix(w io.Writer, client *storage.Client, bucket, prefix, delim string) error {
	ctx := context.Background()
	// [START storage_list_files_with_prefix]
	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimeter, the entire  tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="a/", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="a/" and delim="/", you'll get back:
	//   /a/1.txt
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	// [END storage_list_files_with_prefix]
	return nil
}

func read(client *storage.Client, bucket, object string) ([]byte, error) {
	ctx := context.Background()
	// [START download_file]

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
	// [END download_file]
}

func attrs(client *storage.Client, bucket, object string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	// [START get_metadata]
	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	return attrs, nil
	// [END get_metadata]
}

func makePublic(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START public]
	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	// [END public]
	return nil
}

func move(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START move_file]
	dstName := object + "-rename"

	src := client.Bucket(bucket).Object(object)
	dst := client.Bucket(bucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	if err := src.Delete(ctx); err != nil {
		return err
	}
	// [END move_file]
	return nil
}

func copyToBucket(client *storage.Client, dstBucket, srcBucket, srcObject string) error {
	ctx := context.Background()
	// [START copy_file]
	dstObject := srcObject + "-copy"
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	// [END copy_file]
	return nil
}

func delete(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	// [START delete_file]
	o := client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	// [END delete_file]
	return nil
}

// writeEncryptedObject writes an object encrypted with user-provided AES key to a bucket.
func writeEncryptedObject(client *storage.Client, bucket, object string, secretKey []byte) error {
	ctx := context.Background()

	// [START storage_upload_encrypted_file]
	obj := client.Bucket(bucket).Object(object)
	// Encrypt the object's contents.
	wc := obj.Key(secretKey).NewWriter(ctx)
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END storage_upload_encrypted_file]
	return nil
}

// writeWithKMSKey writes an object encrypted with KMS-provided key to a bucket.
func writeWithKMSKey(client *storage.Client, bucket, object string, keyName string) error {
	ctx := context.Background()

	// [START storage_upload_with_kms_key]
	obj := client.Bucket(bucket).Object(object)
	// Encrypt the object's contents
	wc := obj.NewWriter(ctx)
	wc.KMSKeyName = keyName
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END storage_upload_with_kms_key]
	return nil
}

func readEncryptedObject(client *storage.Client, bucket, object string, secretKey []byte) ([]byte, error) {
	ctx := context.Background()

	// [START storage_download_encrypted_file]
	obj := client.Bucket(bucket).Object(object)
	rc, err := obj.Key(secretKey).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	// [END storage_download_encrypted_file]
	return data, nil
}

func rotateEncryptionKey(client *storage.Client, bucket, object string, key, newKey []byte) error {
	ctx := context.Background()
	// [START storage_rotate_encryption_key]
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	obj := client.Bucket(bucket).Object(object)
	// obj is encrypted with key, we are encrypting it with the newKey.
	_, err = obj.Key(newKey).CopierFrom(obj.Key(key)).Run(ctx)
	if err != nil {
		return err
	}
	// [END storage_rotate_encryption_key]
	return nil
}

func downloadUsingRequesterPays(client *storage.Client, object, bucketName, localpath, billingProjectID string) error {
	ctx := context.Background()
	// [START storage_download_file_requester_pays]
	bucket := client.Bucket(bucketName).UserProject(billingProjectID)
	src := bucket.Object(object)

	f, err := os.OpenFile(localpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	rc, err := src.NewReader(ctx)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, rc); err != nil {
		return err
	}
	if err := rc.Close(); err != nil {
		return err
	}
	fmt.Printf("Downloaded using %v as billing project.\n", billingProjectID)
	// [END storage_download_file_requester_pays]
	return nil
}

// TODO(jbd): Add test for downloadUsingRequesterPays.

const helptext = `usage: ./gcp-bucket-controller -o=bucket:object [subcommand] <args...>
subcommands:
	- write
	- read
	- metadata
	- makepublic
	- delete
`

// GOOGLE_APPLICATION_CREDENTIALS="/Users/jay_liao/go/src/jay-test/in-sports-test-storage.json" ./jay-test -o=jay-test-fat-log:jay/jay.json  read

// GOOGLE_APPLICATION_CREDENTIALS="/Users/jay_liao/go/src/jay-test/in-sports-test-storage.json" ./jay-test -o=jay-test-fat-log:jay.json  write
func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, helptext)
	flag.PrintDefaults()
	os.Exit(2)
}
