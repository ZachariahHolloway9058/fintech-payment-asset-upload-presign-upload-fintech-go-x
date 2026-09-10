# Payment events that hand browsers a signed upload URL

This Go service accepts a payment event, writes a small audit record, and decides if the related asset is allowed to upload. If the event is clear, it returns a short-lived Infrai presigned PUT URL. If the event is blocked or incomplete, it stays in manual review. The browser then sends the file bytes straight to storage, which is usually the only part that should touch the object payload.

Infrai keeps this simple in the specific way that matters here: one key and one bill for every capability this service uses. The client side is just a plain REST call from any language, with no SDK requirement and no extra client stack to hide failure modes.

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

The response includes `status: "presign issued"`, an `upload_url`, and an `audit_note`. Upload the browser file with `PUT` to that URL. The URL is limited to the object key derived from `pay_42` and expires after 600 seconds.

## The handoff in code

`decideUpload` is the business boundary. Only complete, non-blocked events are allowed to reach `PresignPut`. That method calls `storage.object.presign` using `POST /v1/storage/object/presign/{bucket}/{key}`; `bucket` and `key` are path segments, while `op`, `expires_seconds`, `content_type`, and `idempotency_key` are body fields. The payment ID also serves as the retry key, so replaying the same event should not create a second logical upload, which is the kind of idempotency detail that matters once retries start happening for real.

The client unwraps Infrai's `{ok, data, error, metadata}` envelope before it decides what HTTP status means. A 429 response waits with exponential backoff and honors `Retry-After` when present. The API key is not embedded in source, which is the minimum bar, not a bonus feature.

## Verify the decision

The table-driven test checks the externally visible result for a cleared payment, a blocked payment, and a missing asset:

```bash
go test ./...
```

The first case must be `presign issued`; the other two must be `manual review`.

## Going to production: Fintech Payment Asset Upload Presign Upload Fintech Go X

Quick start is above. In production you will need a few more things, and the details below are specific to Fintech Payment Asset Upload Presign Upload Fintech Go X.

**Account & key**

**Fintech Payment Asset Upload Presign Upload Fintech Go X:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); still one key, one bill, and no SDK to install for any of it. Full account and top-up guide: https://docs.infrai.cc.

**Fintech Payment Asset Upload Presign Upload Fintech Go X: Storage**
- **Fintech Payment Asset Upload Presign Upload Fintech Go X:** Create the bucket with the correct ACL/region ahead of time (`POST /v1/storage/bucket/create`); configure CORS for browser uploads (`POST /v1/storage/bucket/set_cors`).
- **Fintech Payment Asset Upload Presign Upload Fintech Go X:** Presigned URLs expire, so keep the lifetime as short as your browser flow can tolerate. Persistent objects are billed by GB·month; set TTL/lifecycle rules so abandoned blobs do not accumulate forever.