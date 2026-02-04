import { useState } from "react";
import { createPet } from "../api";

type AddPetFormProps = {
  onPetCreated: () => Promise<void>;
};

type FormState = {
  name: string;
  species: "CAT" | "DOG" | "FROG";
  ageYears: number;
  pictureUrl: string;
  description: string;
  breederName: string;
  breederEmail: string;
};

const defaultState: FormState = {
  name: "",
  species: "CAT",
  ageYears: 1,
  pictureUrl: "",
  description: "",
  breederName: "",
  breederEmail: "",
};

export default function AddPetForm({ onPetCreated }: AddPetFormProps) {
  const [formState, setFormState] = useState<FormState>(defaultState);
  const [formError, setFormError] = useState<string | null>(null);
  const [formSuccess, setFormSuccess] = useState<string | null>(null);

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault();
    setFormError(null);
    setFormSuccess(null);
    try {
      await createPet({
        ...formState,
        ageYears: Number(formState.ageYears),
      });
      setFormSuccess("Pet created successfully.");
      setFormState(defaultState);
      await onPetCreated();
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message);
      }
    }
  }

  return (
    <section className="form-section">
      <h2>Add a new pet</h2>
      <p className="muted">
        This uses merchant credentials configured in the environment.
      </p>
      <form className="form" onSubmit={handleSubmit}>
        <label>
          Name
          <input
            value={formState.name}
            onChange={(event) =>
              setFormState({ ...formState, name: event.target.value })
            }
            required
          />
        </label>
        <label>
          Species
          <select
            value={formState.species}
            onChange={(event) =>
              setFormState({
                ...formState,
                species: event.target.value as FormState["species"],
              })
            }
          >
            <option value="CAT">Cat</option>
            <option value="DOG">Dog</option>
            <option value="FROG">Frog</option>
          </select>
        </label>
        <label>
          Age (years)
          <input
            type="number"
            min={0}
            value={formState.ageYears}
            onChange={(event) =>
              setFormState({
                ...formState,
                ageYears: Number(event.target.value),
              })
            }
            required
          />
        </label>
        <label>
          Picture URL
          <input
            value={formState.pictureUrl}
            onChange={(event) =>
              setFormState({ ...formState, pictureUrl: event.target.value })
            }
            required
          />
        </label>
        <label>
          Description
          <textarea
            rows={3}
            value={formState.description}
            onChange={(event) =>
              setFormState({ ...formState, description: event.target.value })
            }
            required
          />
        </label>
        <label>
          Breeder name
          <input
            value={formState.breederName}
            onChange={(event) =>
              setFormState({ ...formState, breederName: event.target.value })
            }
            required
          />
        </label>
        <label>
          Breeder email
          <input
            type="email"
            value={formState.breederEmail}
            onChange={(event) =>
              setFormState({ ...formState, breederEmail: event.target.value })
            }
            required
          />
        </label>
        <button className="primary" type="submit">
          Create pet
        </button>
      </form>
      {formError && <div className="status error">{formError}</div>}
      {formSuccess && <div className="status success">{formSuccess}</div>}
    </section>
  );
}
