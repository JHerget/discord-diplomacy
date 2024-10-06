export const OrdersRegex = {
  Move: /^([af])\s*([a-z]{3})\s*-\s*(?:hold|h|([a-z]{3}))$/,
  Support:
    /^([af])\s*([a-z]{3})\s*s\s*([af])\s*([a-z]{3})\s*-\s*(?:hold|h|([a-z]{3}))$/,
  Convoy: /^f\s*([a-z]{3})\s*c\s*a\s*([a-z]{3})\s*-\s*([a-z]{3})$/,
  Reinforce: /^([af])\s*([a-z]{3})$/,
  Disband: /^d\s*([a-z]{3})$/,
};
