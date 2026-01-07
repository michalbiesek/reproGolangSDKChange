# Field Naming Differences: SDK v0.4.0-alpha.13 vs v0.5.0-alpha.1

## Summary

There is a **field naming inconsistency** between what the SDK sends in the REQUEST and what the API returns in the RESPONSE.

## Key Difference

### REQUEST (SDK v0.5.0-alpha.1 sends):
```json
{
  "FileNameSuffix": "`.${C.env[\"CRIBL_WORKER_ID\"]}...`"
}
```
- Uses **PascalCase**: `FileNameSuffix` (capital F, capital N, capital S)

### RESPONSE (API returns):
```json
{
  "fileNameSuffix": "`.${C.env[\"CRIBL_WORKER_ID\"]}...`"
}
```
- Uses **camelCase**: `fileNameSuffix` (lowercase f, lowercase N, capital S)

## Impact

1. **Field Name Mismatch**: The JSON key name differs between request and response
   - Request: `FileNameSuffix` (PascalCase)
   - Response: `fileNameSuffix` (camelCase)

2. **All Other Fields Match**: All other fields appear to have consistent naming between request and response (using camelCase)

3. **Field Values Match**: The actual values are identical, only the key name differs

## Analysis

This suggests:
- The SDK may be serializing Go struct field names directly (PascalCase) for this specific field
- The API consistently returns camelCase field names
- This could indicate a missing or incorrect JSON tag in the SDK model definition for `FileNameSuffix`

## Example Comparison

| Field | REQUEST (SDK) | RESPONSE (API) | Status |
|-------|---------------|----------------|--------|
| `FileNameSuffix` | `FileNameSuffix` | `fileNameSuffix` | ❌ **Mismatch** |
| `baseFileName` | `baseFileName` | `baseFileName` | ✅ Match |
| `bucket` | `bucket` | `bucket` | ✅ Match |
| `compress` | `compress` | `compress` | ✅ Match |
| All other fields | camelCase | camelCase | ✅ Match |

## Recommendation

This inconsistency should be investigated in the SDK to ensure:
1. Request serialization uses consistent camelCase naming
2. The `FileNameSuffix` field has the correct JSON tag: `json:"fileNameSuffix"`
3. SDK version v0.5.0-alpha.1 properly handles field name serialization

