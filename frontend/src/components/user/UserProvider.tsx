let mock;

if   (import.meta.env.DEV) { mock = "MockUserProvider"; }
else { mock = "ApiUserProvider"; }

const { UserProvider } = await import(`./providers/${mock}.tsx`);
export { UserProvider };
