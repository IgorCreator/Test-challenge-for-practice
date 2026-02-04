import { useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { fetchPurchasedPets, fetchStorePets, purchasePets } from "../api";
import { Pet, PurchaseError } from "../types";
import AddPetForm from "../components/AddPetForm";
import PetGrid from "../components/PetGrid";
import PurchaseHistory from "../components/PurchaseHistory";
import StoreHeader from "../components/StoreHeader";
import Tabs, { TabKey } from "../components/Tabs";

export default function App() {
  const { slug } = useParams();
  const [pets, setPets] = useState<Pet[]>([]);
  const [purchased, setPurchased] = useState<Pet[]>([]);
  const [cart, setCart] = useState<Record<string, Pet>>({});
  const [activeTab, setActiveTab] = useState<TabKey>("store");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [checkoutErrors, setCheckoutErrors] = useState<PurchaseError[]>([]);

  const cartItems = useMemo(() => Object.values(cart), [cart]);

  useEffect(() => {
    let mounted = true;
    setLoading(true);
    Promise.all([fetchStorePets(slug), fetchPurchasedPets(slug)])
      .then(([available, purchasedList]) => {
        if (!mounted) return;
        setPets(available);
        setPurchased(purchasedList);
        setError(null);
      })
      .catch((err) => {
        if (!mounted) return;
        setError(err.message);
      })
      .finally(() => {
        if (!mounted) return;
        setLoading(false);
      });
    return () => {
      mounted = false;
    };
  }, [slug]);

  function toggleCart(pet: Pet) {
    setCart((prev) => {
      const next = { ...prev };
      if (next[pet.id]) {
        delete next[pet.id];
      } else {
        next[pet.id] = pet;
      }
      return next;
    });
  }

  async function handleCheckout() {
    setCheckoutErrors([]);
    try {
      const result = await purchasePets(cartItems.map((p) => p.id), slug);
      if (result.errors.length > 0) {
        setCheckoutErrors(result.errors);
      } else {
        setCart({});
      }
      const refreshed = await fetchStorePets(slug);
      setPets(refreshed);
      const purchasedList = await fetchPurchasedPets(slug);
      setPurchased(purchasedList);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setCheckoutErrors([{ petName: "Checkout", message: err.message }]);
      }
    }
  }

  async function refreshPets() {
    const refreshed = await fetchStorePets(slug);
    setPets(refreshed);
  }

  return (
    <div className="page">
      <StoreHeader
        cartCount={cartItems.length}
        checkoutDisabled={cartItems.length === 0}
        onCheckout={handleCheckout}
      />

      <Tabs active={activeTab} onChange={setActiveTab} />

      {loading && <div className="status">Loading pets...</div>}
      {error && <div className="status error">{error}</div>}

      {checkoutErrors.length > 0 && activeTab === "store" && (
        <div className="alert">
          <strong>Some pets are no longer available:</strong>
          <ul>
            {checkoutErrors.map((err) => (
              <li key={`${err.petName}-${err.message}`}>
                {err.petName}: {err.message}
              </li>
            ))}
          </ul>
        </div>
      )}

      {activeTab === "store" && (
        <PetGrid pets={pets} cart={cart} onToggleCart={toggleCart} />
      )}

      {activeTab === "add" && <AddPetForm onPetCreated={refreshPets} />}

      {activeTab === "history" && <PurchaseHistory pets={purchased} />}
    </div>
  );
}
