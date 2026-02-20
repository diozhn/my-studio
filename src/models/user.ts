export interface User {
  id: number;
  username: string;
  password: string;
  email: string;
  superuser: boolean;
  google_id: string | null;
  instagram_id: string | null;
  twitter_id: string | null;
  refresh_token: string | null;
  created_at: Date;
}


