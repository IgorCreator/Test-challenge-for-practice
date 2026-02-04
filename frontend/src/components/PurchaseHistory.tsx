import { Pet } from "../types";

type PurchaseHistoryProps = {
  pets: Pet[];
};

export default function PurchaseHistory({ pets }: PurchaseHistoryProps) {
  return (
    <section className="history">
      <h2>Purchased history</h2>
      {pets.length === 0 ? (
        <p className="muted">No purchases yet.</p>
      ) : (
        <div className="history-grid">
          {pets.map((pet) => (
            <div key={pet.id} className="history-card">
              <div>
                <strong>{pet.name}</strong> ({pet.species})
              </div>
              <div className="muted">
                Purchased:{" "}
                {pet.purchasedAt
                  ? new Date(pet.purchasedAt).toLocaleString()
                  : "—"}
              </div>
              <div className="code">ID: {pet.id}</div>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}
