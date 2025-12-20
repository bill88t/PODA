

let mock;

if   (import.meta.env.DEV) { mock = "MockUserProvider"; }
else { mock = "MockUserProvider"; }

const { UserProvider } = await import(`./providers/${mock}.tsx`);
export { UserProvider };
