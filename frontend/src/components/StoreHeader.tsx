type StoreHeaderProps = {
  cartCount: number;
  onCheckout: () => void;
  checkoutDisabled: boolean;
};

const initialMessage =
  "Browse pets and add them to your cart. Checkout will purchase everything at once.";

export default function StoreHeader({
  cartCount,
  onCheckout,
  checkoutDisabled,
}: StoreHeaderProps) {
  return (
    <header className="header">
      <div>
        <h1>Nimble Pet Store</h1>
        <p className="subtitle">{initialMessage}</p>
      </div>
      <div className="cart">
        <div className="cart-count">{cartCount} in cart</div>
        <button className="primary" disabled={checkoutDisabled} onClick={onCheckout}>
          Checkout
        </button>
      </div>
    </header>
  );
}
