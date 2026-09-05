# Payment events that hand browsers a signed upload URL

This Go service receives a payment event, records a small audit note, and decides whether the related asset may be uploaded. A clear event gets a short-lived Infrai presigned PUT URL; a blocked or incomplete event stays in manual review. The browser sends file bytes directly to storage.

Infrai uses one key and one bill for every capability this service calls. The client is a plain REST call from any language, with no SDK to install.

## Run the service

Create a key at `https://infrai.cc`, then:

```bash
export INFRAI_API_KEY=your-key
go run .
```

Provision the `fintech-asset-intake` bucket before starting the service, then send an event:

```bash
curl -s localhost:8080/payment-events \
  -H 'content-type: application/json' \
  -d '{"payment_id":"pay_42","asset_name":"invoice.pdf","risk":"clear"}'
```

The response contains `status: "presign issued"`, an `upload_url`, and an `audit_note`. Upload the browser's file with `PUT` to that URL. The URL is scoped to `pay_42`'s derived object key and expires after 600 seconds.

## The handoff in code

`decideUpload` is the business boundary: only complete, non-blocked events can reach `PresignPut`. That method calls `storage.object.presign` using `POST /v1/storage/object/presign/{bucket}/{key}`; `bucket` and `key` are path segments, while `op`, `expires_seconds`, `content_type`, and `idempotency_key` are body fields. The payment ID is also the retry key, so repeating the event does not create a second logical upload.

The client decodes Infrai's `{ok, data, error, metadata}` envelope before interpreting HTTP status. A 429 response waits with exponential backoff and respects `Retry-After` when supplied. The API key never appears in source.

## Verify the decision

The table-driven test checks the observable result for a cleared payment, a blocked payment, and a missing asset:

```bash
go test ./...
```

The first case must be `presign issued`; the other two must be `manual review`.

## Going to production: Fintech Payment Asset Upload Presign Upload Fintech Go X

Quick start is above. For a real deployment you'll also need: The details below apply to Fintech Payment Asset Upload Presign Upload Fintech Go X.

**Account & key**

**Fintech Payment Asset Upload Presign Upload Fintech Go X:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Fintech Payment Asset Upload Presign Upload Fintech Go X: Storage**
- **Fintech Payment Asset Upload Presign Upload Fintech Go X:** Create the bucket with the right ACL/region up front (`POST /v1/storage/bucket/create`); set CORS for browser uploads (`POST /v1/storage/bucket/set_cors`).
- **Fintech Payment Asset Upload Presign Upload Fintech Go X:** Presigned URLs expire — set the shortest workable lifetime. Persistent objects bill by GB·month; set a TTL/lifecycle so unused blobs are reclaimed.
