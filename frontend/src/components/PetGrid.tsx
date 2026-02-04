import { Pet } from "../types";

type PetGridProps = {
  pets: Pet[];
  cart: Record<string, Pet>;
  onToggleCart: (pet: Pet) => void;
};

export default function PetGrid({ pets, cart, onToggleCart }: PetGridProps) {
  return (
    <div className="grid">
      {pets.map((pet) => {
        const inCart = Boolean(cart[pet.id]);
        return (
          <div key={pet.id} className="card">
            <img src={pet.pictureUrl} alt={pet.name} />
            <div className="card-body">
              <div className="card-title">
                <h3>{pet.name}</h3>
                <span className="tag">{pet.species}</span>
              </div>
              <p className="muted">{pet.description}</p>
              <div className="meta">
                <span>Age: {pet.ageYears} years</span>
                <span>Breeder: {pet.breederName}</span>
                <span className="code">ID: {pet.id}</span>
              </div>
            </div>
            <button
              className={inCart ? "secondary" : "primary"}
              onClick={() => onToggleCart(pet)}
            >
              {inCart ? "Remove from cart" : "Add to cart"}
            </button>
          </div>
        );
      })}
    </div>
  );
}
