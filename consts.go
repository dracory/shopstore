package shopstore

const CATEGORY_STATUS_ACTIVE = "active"
const CATEGORY_STATUS_DRAFT = "draft"
const CATEGORY_STATUS_INACTIVE = "inactive"

const COLUMN_AMOUNT = "amount"
const COLUMN_CODE = "code"
const COLUMN_CREATED_AT = "created_at"
const COLUMN_CUSTOMER_ID = "customer_id"
const COLUMN_DESCRIPTION = "description"
const COLUMN_ENDS_AT = "ends_at"
const COLUMN_ENTITY_ID = "entity_id"
const COLUMN_ID = "id"
const COLUMN_MEDIA_TYPE = "media_type"
const COLUMN_MEDIA_URL = "media_url"
const COLUMN_MEMO = "memo"
const COLUMN_METAS = "metas"
const COLUMN_ORDER_ID = "order_id"
const COLUMN_PARENT_ID = "parent_id"
const COLUMN_PRICE = "price"
const COLUMN_PRODUCT_ID = "product_id"
const COLUMN_QUANTITY = "quantity"
const COLUMN_SEQUENCE = "sequence"
const COLUMN_SOFT_DELETED_AT = "soft_deleted_at"
const COLUMN_SHORT_DESCRIPTION = "short_description"
const COLUMN_STARTS_AT = "starts_at"
const COLUMN_STATUS = "status"
const COLUMN_TYPE = "type"
const COLUMN_TITLE = "title"
const COLUMN_UPDATED_AT = "updated_at"

const MEDIA_STATUS_DRAFT = "draft"
const MEDIA_STATUS_ACTIVE = "active"
const MEDIA_STATUS_INACTIVE = "inactive"

const MEDIA_TYPE_IMAGE_JPG = "image/jpeg"
const MEDIA_TYPE_IMAGE_PNG = "image/png"
const MEDIA_TYPE_VIDEO_MP4 = "video/mp4"

// Customer has completed the checkout process, but payment has yet to be confirmed.
const ORDER_STATUS_AWAITING_PAYMENT = "awaiting_payment"

// Customer has completed the checkout process and payment has been confirmed.
const ORDER_STATUS_AWAITING_FULFILLMENT = "awaiting_fulfillment"

// Order has been packaged and is awaiting customer pickup from a seller-specified location.
const ORDER_STATUS_AWAITING_PICKUP = "awaiting_pickup"

// Order has been pulled and packaged and is awaiting collection from a shipping provider.
const ORDER_STATUS_AWAITING_SHIPMENT = "awaiting_shipment"

// Seller has cancelled an order, due to a stock inconsistency or other reasons. Cancelling an order will not refund the order.
const ORDER_STATUS_CANCELLED = "cancelled"

// Order has been shipped/picked up, and receipt is confirmed; client has paid for their digital product, and their file(s) are available for download.
const ORDER_STATUS_COMPLETED = "completed"

// Seller has marked the order as declined.
const ORDER_STATUS_DECLINED = "declined"

// Customer has initiated a dispute resolution process for the PayPal transaction that paid for the order or the seller has marked the order as a fraudulent order.
const ORDER_STATUS_DISPUTED = "disputed"

// Order on hold while some aspect, such as tax-exempt documentation, is manually confirmed.
// Orders with this status must be updated manually.
const ORDER_STATUS_MANUAL_VERIFICATION_REQUIRED = "manual_verification_required"

// Seller has partially refunded the order.
const ORDER_STATUS_PARTIALLY_REFUNDED = "partially_refunded"

// Only some items in the order have been shipped.
const ORDER_STATUS_PARTIALLY_SHIPPED = "partially_shipped"

// Customer started the checkout process but did not complete it.
const ORDER_STATUS_PENDING = "pending"

// Seller has used the Refund action to refund the whole order. A listing of all orders with a "Refunded" status can be found under the More tab of the View Orders screen.
const ORDER_STATUS_REFUNDED = "refunded"

// Order has been shipped, but receipt has not been confirmed; seller has used the Ship Items action. A listing of all orders with a "Shipped" status can be found under the More tab of the View Orders screen.
const ORDER_STATUS_SHIPPED = "shipped"

const COLUMN_VARIANT_MATRIX_SCHEMA = "variant_matrix_schema"

const COLUMN_VARIANT_MATRIX_VALUES = "variant_matrix_values"

const PRODUCT_STATUS_DRAFT = "draft"

const PRODUCT_STATUS_ACTIVE = "active"

const PRODUCT_STATUS_DISABLED = "disabled"

const PRODUCT_STATUS_PARENT = "parent"
