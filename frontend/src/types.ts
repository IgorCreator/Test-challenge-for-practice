export type Pet = {
  id: string;
  name: string;
  species: "CAT" | "DOG" | "FROG";
  ageYears: number;
  pictureUrl: string;
  description: string;
  breederName: string;
  breederEmail: string;
  createdAt: string;
  purchasedAt?: string | null;
};

export type PurchaseError = {
  petName: string;
  message: string;
};
