import { Pet, PurchaseError } from "./types";

const API_URL = import.meta.env.VITE_API_URL as string;
const CUSTOMER_USER = import.meta.env.VITE_CUSTOMER_USER as string;
const CUSTOMER_PASS = import.meta.env.VITE_CUSTOMER_PASS as string;
const MERCHANT_USER = import.meta.env.VITE_MERCHANT_USER as string;
const MERCHANT_PASS = import.meta.env.VITE_MERCHANT_PASS as string;
const STORE_SLUG = import.meta.env.VITE_STORE_SLUG as string;

type GraphQLResponse<T> = {
  data?: T;
  errors?: { message: string }[];
};

function basicAuthHeader(username: string, password: string) {
  const token = btoa(`${username}:${password}`);
  return `Basic ${token}`;
}

async function request<T>(
  query: string,
  variables: Record<string, unknown>,
  auth?: { username: string; password: string }
): Promise<T> {
  const credentials = auth ?? { username: CUSTOMER_USER, password: CUSTOMER_PASS };
  const res = await fetch(API_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: basicAuthHeader(credentials.username, credentials.password),
    },
    body: JSON.stringify({ query, variables }),
  });
  if (!res.ok) {
    throw new Error(`Request failed (${res.status})`);
  }
  const payload = (await res.json()) as GraphQLResponse<T>;
  if (payload.errors && payload.errors.length > 0) {
    throw new Error(payload.errors.map((e) => e.message).join(", "));
  }
  if (!payload.data) {
    throw new Error("No data returned");
  }
  return payload.data;
}

export async function fetchStorePets(slug = STORE_SLUG): Promise<Pet[]> {
  const query = `
    query StorePets($storeSlug: String!) {
      storePets(storeSlug: $storeSlug) {
        id
        name
        species
        ageYears
        pictureUrl
        description
        breederName
        breederEmail
        createdAt
        purchasedAt
      }
    }
  `;
  const data = await request<{ storePets: Pet[] }>(query, { storeSlug: slug });
  return data.storePets;
}

export async function fetchPurchasedPets(slug = STORE_SLUG): Promise<Pet[]> {
  const query = `
    query PurchasedPets($storeSlug: String!) {
      purchasedPets(storeSlug: $storeSlug) {
        id
        name
        species
        ageYears
        pictureUrl
        description
        breederName
        createdAt
        purchasedAt
      }
    }
  `;
  const data = await request<{ purchasedPets: Pet[] }>(query, { storeSlug: slug });
  return data.purchasedPets;
}

export async function purchasePets(
  petIds: string[],
  slug = STORE_SLUG
): Promise<{ purchasedIds: string[]; errors: PurchaseError[] }> {
  const query = `
    mutation PurchasePets($input: PurchasePetsInput!) {
      purchasePets(input: $input) {
        purchasedIds
        errors {
          petName
          message
        }
      }
    }
  `;
  const data = await request<{
    purchasePets: { purchasedIds: string[]; errors: PurchaseError[] };
  }>(query, { input: { storeSlug: slug, petIds } });
  return data.purchasePets;
}

export async function createPet(input: {
  name: string;
  species: "CAT" | "DOG" | "FROG";
  ageYears: number;
  pictureUrl: string;
  description: string;
  breederName: string;
  breederEmail: string;
}): Promise<Pet> {
  const query = `
    mutation CreatePet($input: CreatePetInput!) {
      createPet(input: $input) {
        id
        name
        species
        ageYears
        pictureUrl
        description
        breederName
        breederEmail
        createdAt
      }
    }
  `;
  const data = await request<{ createPet: Pet }>(
    query,
    { input },
    { username: MERCHANT_USER, password: MERCHANT_PASS }
  );
  return data.createPet;
}
